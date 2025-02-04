package handlers

import (
	"ai-sdk-server/config"
	"ai-sdk-server/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/postgrest-go"
	"net/http"
	"os"
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

//func (h *RequestHandler) CreateModelRequest(c *fiber.Ctx) error {
//	var input struct {
//		ModelType string                 `json:"model_type"`
//		InputData map[string]interface{} `json:"input_data"`
//	}
//
//	if err := c.BodyParser(&input); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"error": "Invalid request body",
//		})
//	}
//
//	userID := c.Locals("user_id").(string)
//	requestID := uuid.New()
//
//	request := models.ModelRequest{
//		ID:        requestID,
//		UserID:    uuid.MustParse(userID),
//		Status:    "pending",
//		InputData: input.InputData,
//		ModelType: input.ModelType,
//		CreatedAt: time.Now(),
//	}
//
//	_, _, err := h.dbClient.From("model_requests").
//		Insert(request, false, "", "representation", "exact").
//		Execute()
//
//	if err != nil {
//		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
//			"error": "request not initialized in database",
//		})
//	}
//
//	edgeFunctionPayload := map[string]interface{}{
//		"id":      requestID.String(),
//		"input":   input.InputData,
//		"modelId": input.ModelType, // The Edge Function expects modelId
//	}
//
//	httpClient := &http.Client{
//		Timeout: time.Second * 30, // Adjust timeout as needed
//	}
//
//	payloadBytes, err := json.Marshal(edgeFunctionPayload)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to prepare edge function request",
//		})
//	}
//
//	edgeFunctionURL := "https://tpuhjnicfmhvgoufjvvn.supabase.co/functions/v1/huggingface-models"
//	req, err := http.NewRequest("POST", edgeFunctionURL, bytes.NewBuffer(payloadBytes))
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to create edge function request",
//		})
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
//
//	resp, err := httpClient.Do(req)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Edge function request failed",
//		})
//	}
//	defer resp.Body.Close()
//
//	var edgeFunctionResponse map[string]interface{}
//	if err := json.NewDecoder(resp.Body).Decode(&edgeFunctionResponse); err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to parse edge function response",
//		})
//	}
//
//	var updatedRequest models.ModelRequest
//
//	_, _, err = h.dbClient.From("model_requests").
//		Select("*", "exact", false).
//		Eq("id", requestID.String()).
//		Execute()
//
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to fetch updated request",
//		})
//	}
//
//	return c.Status(fiber.StatusOK).JSON(updatedRequest)
//}

func (h *RequestHandler) ListRequests(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var requests []models.ModelRequest
	_, _, err := h.dbClient.From("model_requests").
		Select("*", "exact", false).
		Eq("user_id", userID).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch requests",
		})
	}

	return c.JSON(requests)
}

//func (h *RequestHandler) CreateModelRequest(c *fiber.Ctx) error {
//	var input struct {
//		ModelType string                 `json:"model_type"`
//		InputData map[string]interface{} `json:"input_data"`
//	}
//
//	if err := c.BodyParser(&input); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"error": "Invalid request body",
//		})
//	}
//
//	userID := c.Locals("user_id").(string)
//	requestID := uuid.New()
//
//	request := models.ModelRequest{
//		ID:        requestID,
//		UserID:    uuid.MustParse(userID),
//		Status:    "pending",
//		InputData: input.InputData,
//		ModelType: input.ModelType,
//		CreatedAt: time.Now(),
//		// Initialize other fields to prevent null issues
//		OutputData:     make(map[string]interface{}),
//		ErrorMsg:       "",
//		TokenUsed:      0,
//		TokenCount:     0,
//		ProcessingTime: 0,
//		CompletedAT:    time.Now(),
//	}
//
//	resp, _, err := h.dbClient.From("model_requests").
//		Insert(request, false, "", "representation", "exact").
//		Single(). // We want a single row returned
//		Execute()
//
//	if err != nil {
//		// Log the error details for debugging
//		fmt.Printf("Database error: %v\n", err)
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to create database record: " + err.Error(),
//		})
//	}
//
//	// Decode the response to verify the record was created successfully
//	var createdRequest models.ModelRequest
//	if err := json.NewDecoder(bytes.NewBuffer(resp)).Decode(&createdRequest); err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to parse database response",
//		})
//	}
//
//	// Verify that we got back a valid record with an ID
//	if createdRequest.ID == uuid.Nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to create request - no ID returned",
//		})
//	}
//
//	edgeFunctionPayload := map[string]interface{}{
//		"id":      requestID.String(),
//		"input":   input.InputData,
//		"modelId": input.ModelType, // The Edge Function expects modelId
//	}
//
//	httpClient := &http.Client{
//		Timeout: time.Second * 30, // Adjust timeout as needed
//	}
//
//	payloadBytes, err := json.Marshal(edgeFunctionPayload)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to prepare edge function request",
//		})
//	}
//
//	edgeFunctionURL := "https://tpuhjnicfmhvgoufjvvn.supabase.co/functions/v1/hf-inference"
//	req, err := http.NewRequest("POST", edgeFunctionURL, bytes.NewBuffer(payloadBytes))
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to create edge function request",
//		})
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
//
//	response, err := httpClient.Do(req)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Edge function request failed",
//		})
//	}
//	defer response.Body.Close()
//
//	var edgeFunctionResponse map[string]interface{}
//	if err := json.NewDecoder(response.Body).Decode(&edgeFunctionResponse); err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to parse edge function response",
//		})
//	}
//
//	var updatedRequest models.ModelRequest
//
//	_, _, err = h.dbClient.From("model_requests").
//		Select("*", "exact", false).
//		Eq("id", requestID.String()).
//		Execute()
//
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": "Failed to fetch updated request",
//		})
//	}
//
//	return c.Status(fiber.StatusOK).JSON(updatedRequest)
//}

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
		// Initialize other fields to prevent null issues
		OutputData:     make(map[string]interface{}),
		ErrorMsg:       "",
		TokenUsed:      0,
		TokenCount:     0,
		ProcessingTime: 0,
		CompletedAT:    time.Now(),
	}

	resp, _, err := h.dbClient.From("model_requests").
		Insert(request, false, "", "representation", "exact").
		Single(). // We want a single row returned
		Execute()

	if err != nil {
		// Log the error details for debugging
		fmt.Printf("Database error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create database record: " + err.Error(),
		})
	}

	// Decode the response to verify the record was created successfully
	var createdRequest models.ModelRequest
	if err := json.NewDecoder(bytes.NewBuffer(resp)).Decode(&createdRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse database response",
		})
	}

	// Verify that we got back a valid record with an ID
	if createdRequest.ID == uuid.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request - no ID returned",
		})
	}

	edgeFunctionPayload := map[string]interface{}{
		"id":      requestID.String(),
		"input":   input.InputData,
		"modelId": input.ModelType, // The Edge Function expects modelId
	}

	httpClient := &http.Client{
		Timeout: time.Second * 360, // Adjust timeout as needed
	}

	payloadBytes, err := json.Marshal(edgeFunctionPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to prepare edge function request",
		})
	}

	edgeFunctionURL := "https://tpuhjnicfmhvgoufjvvn.supabase.co/functions/v1/hf-inference"
	req, err := http.NewRequest("POST", edgeFunctionURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create edge function request",
		})
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

	response, err := httpClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Edge function request failed",
		})
	}
	defer response.Body.Close()

	var edgeFunctionResponse map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&edgeFunctionResponse); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse edge function response",
		})
	}

	var updatedRequest models.ModelRequest

	_, _, err = h.dbClient.From("model_requests").
		Select("*", "exact", false).
		Eq("id", requestID.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch updated request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedRequest)
}
