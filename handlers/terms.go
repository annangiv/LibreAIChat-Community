// handlers/terms.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// TermsHandler renders the terms page
func TermsHandler(c *fiber.Ctx) error {
	return c.Render("terms", fiber.Map{
		"Title": "Terms",
	}, "layouts/main")
}
