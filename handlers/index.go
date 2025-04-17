// handlers/index.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// IndexHandler renders the home page
func IndexHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Home",
	}, "layouts/main")
}
