// handlers/privacy.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// PrivacyHandler renders the privacy page
func PrivacyHandler(c *fiber.Ctx) error {
	return c.Render("privacy", fiber.Map{
		"Title": "Privacy",
	}, "layouts/main")
}
