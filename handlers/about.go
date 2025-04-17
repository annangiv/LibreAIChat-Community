// handlers/about.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// AboutHandler renders the about page
func AboutHandler(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title": "About",
	}, "layouts/main")
}
