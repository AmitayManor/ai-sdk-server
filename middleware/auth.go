package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
)

//func AuthRequired() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//
//		authHeader := c.Get("Authorization")
//		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
//			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//				"error": "Missing or invalid authorization header",
//			})
//		}
//
//		token := strings.TrimPrefix(authHeader, "Bearer ")
//
//		client := config.GetSupabaseClient()
//		client.Auth.WithToken(token)
//		user, err := client.Auth.GetUser()
//
//		if err != nil {
//			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//				"error": "Invalid token",
//			})
//		}
//
//		c.Locals("user_id", user.ID)
//
//		return c.Next()
//	}
//}

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check authorization header exists and has correct format
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid authorization header",
			})
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the JWT token to validate it and get claims
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key used to sign the token
			// This should match the JWT secret used by Supabase Auth
			return []byte(os.Getenv("SUPABASE_JWT_SECRET")), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Get the sub claim which contains the user ID
		sub, ok := claims["sub"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Store user ID in context
		c.Locals("user_id", sub)

		return c.Next()
	}
}
