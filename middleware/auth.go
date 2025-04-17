package middleware

import (
	"LibreAI/internal/authentication"
	"LibreAI/internal/services"

	"github.com/gofiber/fiber/v2"
)

func RequireAuth(c *fiber.Ctx) error {
	if authentication.GetUserID(c) == 0 {
		// Not logged in
		return c.Redirect("/auth")
	}
	return c.Next()
}

func RequireAdmin(c *fiber.Ctx) error {
	userID := authentication.GetUserID(c)
	if userID == 0 {
		return c.Redirect("/login")
	}

	user, err := services.GetUserByID(userID)
	if err != nil || user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).SendString("Forbidden")
	}

	return c.Next()
}
