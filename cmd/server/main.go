package main

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"LibreAI/internal/database"
	"LibreAI/models"
	"LibreAI/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2" // Add this import
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

var (
	ModelCache = cache.New(1*time.Hour, 10*time.Minute)
)

func main() {
	_ = godotenv.Load()

	if maxProcsEnv := os.Getenv("MAX_PROCS"); maxProcsEnv != "" {
		if maxProcs, err := strconv.Atoi(maxProcsEnv); err == nil {
			runtime.GOMAXPROCS(maxProcs)
		}
	}

	db := database.Get()
	db.AutoMigrate(&models.User{}, &models.UserStats{}, &models.Model{}, &models.Question{})

	// Initialize the template engine
	viewsEngine := html.New("./views", ".html")

	// In development, this prevents caching of templates
	// Remove in production for better performance
	viewsEngine.Reload(true)

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Ollama UI",
		Views:         viewsEngine, // Add the template engine to the Fiber config
	})

	// Add static file handling
	app.Static("/static", "./static")

	routers.SetupRoutes(app, db)

	log.Printf("Using up to %d CPU cores\n", runtime.GOMAXPROCS(0))
	log.Println("Listening on :3000")
	log.Fatal(app.Listen(":3000"))
}
