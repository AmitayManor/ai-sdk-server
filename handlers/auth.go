package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

type AuthHandler struct {
	supabaseClient *supabase.Client
}

func NewAuthHandler(supabaseClient *supabase.Client) *AuthHandler {
	return &AuthHandler{
		supabaseClient: supabaseClient,
	}
}
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {

	var input types.SignupRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid body request",
		})
	}

	_, err := h.supabaseClient.Auth.Signup(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "signup failed",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Signup successful. Please check your email for verification.",
	})
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	authResponse, err := h.supabaseClient.Auth.SignInWithEmailPassword(input.Email, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid credentials",
			"details": err.Error(),
		})
	}

	// Return the access token
	return c.JSON(fiber.Map{
		"token": authResponse.AccessToken,
	})
}
