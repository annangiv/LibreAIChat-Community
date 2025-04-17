package authentication

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Store user ID in session cookie or context (you can customize this)
func AuthStore(c *fiber.Ctx, userID int) {
	c.Cookie(&fiber.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(userID),
	})
}

func AuthClear(c *fiber.Ctx) {
	c.ClearCookie("user_id")
}

func GetUserID(c *fiber.Ctx) int {
	idStr := c.Cookies("user_id")
	id, _ := strconv.Atoi(idStr)
	return id
}
