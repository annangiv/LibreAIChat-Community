package routers

import (
	"LibreAI/authentication"
	"LibreAI/handlers"
	"LibreAI/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Public pages
	app.Get("/", handlers.IndexHandler)
	app.Get("/about", handlers.AboutHandler)
	app.Get("/terms", handlers.TermsHandler)
	app.Get("/privacy", handlers.PrivacyHandler)
	app.Get("/community", handlers.CommunityHandler)
	app.Get("/goodbye", handlers.GoodByeHandler)

	// Minimal auth routes (GitHub only for simplicity)
	handlers.InitOAuth()
	app.Get("/auth", handlers.AuthPageHandler)
	app.Get("/auth/github", handlers.GitHubLogin)
	app.Get("/auth/github/callback", handlers.GitHubCallback)
	app.Post("/auth/cancel", handlers.CancelOAuth)
	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("user_id")
		return c.Redirect("/")
	})

	// Authenticated routes
	app.Get("/me", handlers.MeHandler)
	app.Get("/account", middleware.RequireAuth, handlers.AccountHandler)
	app.Get("/ask", middleware.RequireAuth, handlers.AskPageHandler)
	app.Post("/ask", middleware.RequireAuth, handlers.AskHandler)
	app.Post("/process-consent", handlers.ProcessConsent)
	app.Post("/delete-account", middleware.RequireAuth, handlers.DeleteAccount)

	// WebSocket for streaming
	app.Use("/ws/ollama", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			userID := authentication.GetUserID(c)
			if userID == 0 {
				return fiber.ErrUnauthorized
			}
			c.Locals("user_id", userID)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/ollama", websocket.New(handlers.HandlePromptStream))

	// Error handling
	app.Use(func(c *fiber.Ctx) error {
		return handlers.ErrorHandler(c)
	})
}
