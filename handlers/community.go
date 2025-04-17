// handlers/community.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// CommunityHandler renders the community page
func CommunityHandler(c *fiber.Ctx) error {
	return c.Render("community", fiber.Map{
		"Title": "Community",
	}, "layouts/main")
}
