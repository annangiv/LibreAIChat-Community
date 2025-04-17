package authentication

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthStore(c *fiber.Ctx, userID int) {
	c.Cookie(&fiber.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(userID),
		Path:     "/",
		HTTPOnly: true,
		Secure:   false, // set to true in production over HTTPS
		SameSite: "Lax", // ensures it survives redirects like OAuth
	})
}

func AuthLogout(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "user_id",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   false, // match your AuthStore setting
		SameSite: "Lax",
		Expires:  time.Now().Add(-time.Hour), // expire immediately
	})
}

func GetUserID(c *fiber.Ctx) int {
	cookie := c.Cookies("user_id") // Updated method
	if cookie == "" {
		return 0
	}

	userID, err := strconv.Atoi(cookie)
	if err != nil {
		return 0
	}

	return userID
}
