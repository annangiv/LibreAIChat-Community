// handlers/goodbye.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// AboutHandler renders the goodbye page
func GoodByeHandler(c *fiber.Ctx) error {
	return c.Render("goodbye", fiber.Map{
		"Title": "Goodbye",
	}, "layouts/main")
}
