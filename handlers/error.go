// handlers/error.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// AboutHandler renders the about page
func ErrorHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).Render("error", fiber.Map{
		"Title":   "Page Not Found",
		"Code":    404,
		"Message": "The page you're looking for doesn't exist.",
	}, "layouts/main")
}
