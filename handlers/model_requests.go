package handlers

import (
	"ai-sdk-server/config"
	"ai-sdk-server/models"
	"bytes"
	_ "context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/postgrest-go"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type RequestHandler struct {
	dbClient *postgrest.Client
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		dbClient: config.GetDBClient(),
	}
}

func (h *RequestHandler) CreateModelRequest(c *fiber.Ctx) error {
	var input struct {
		ModelType string                 `json:"model_type"`
		InputData map[string]interface{} `json:"input_data"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID := c.Locals("user_id").(string)
	requestID := uuid.New()

	request := models.ModelRequest{
		ID:        requestID,
		UserID:    uuid.MustParse(userID),
		Status:    "pending",
		InputData: input.InputData,
		ModelType: input.ModelType,
		CreatedAt: time.Now(),
	}

	_, _, err := h.dbClient.From("model_requests").
		Insert(request, false, "", "representation", "exact").
		Execute()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Request not initialized in database",
		})
	}

	prompt, ok := input.InputData["prompt"].(string)
	if !ok || prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input: 'prompt' is required as a non-empty string",
		})
	}

	edgeFunctionPayload := map[string]interface{}{
		"id":        requestID.String(),
		"input":     prompt,
		"modelType": input.ModelType,
		"userId":    userID,
	}

	payloadBytes, err := json.Marshal(edgeFunctionPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to prepare edge function request",
		})
	}

	edgeFunctionURL := "https://tpuhjnicfmhvgoufjvvn.supabase.co/functions/v1/hf-inference"
	httpClient := &http.Client{
		Timeout: 90 * time.Second,
	}

	req, err := http.NewRequest("POST", edgeFunctionURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create edge function request",
		})
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

	resp, err := httpClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			log.Printf("Edge function HTTP call timed out; proceeding to poll DB for updated status.")
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  "Edge function request failed",
				"detail": err.Error(),
			})
		}
	} else {
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  "Edge function request failed",
				"detail": string(bodyBytes),
			})
		}
		resp.Body.Close()
	}

	var updatedRequest models.ModelRequest
	timeout := 180 * time.Second
	interval := 2 * time.Second
	startTime := time.Now()

	for time.Since(startTime) < timeout {
		res, _, err := h.dbClient.From("model_requests").
			Select("*", "exact", false).
			Eq("id", requestID.String()).
			Execute()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch updated request",
			})
		}

		var records []models.ModelRequest
		if err := json.Unmarshal(res, &records); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse updated request",
			})
		}

		if len(records) > 0 {
			updatedRequest = records[0]
			if updatedRequest.Status != "pending" {
				break
			}
		}
		time.Sleep(interval)
	}

	if updatedRequest.Status == "pending" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Processing still in progress. Please poll again after some time.",
			"record":  updatedRequest,
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedRequest)

}

func (h *RequestHandler) ListRequests(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	res, _, err := h.dbClient.From("model_requests").
		Select("*", "exact", false).
		Eq("user_id", userID).
		Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch requests",
		})
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(res, &records); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Failed to parse requests",
			"detail": err.Error(),
			"raw":    string(res),
		})
	}

	return c.JSON(records)
}
