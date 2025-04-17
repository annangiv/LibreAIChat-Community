package handlers

import (
	"LibreAI/internal/authentication"
	"LibreAI/internal/database"
	"LibreAI/internal/services"
	"LibreAI/models"
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type GenerateRequest struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Stream bool     `json:"stream"`
	Stop   []string `json:"stop,omitempty"`
}

type OllamaStreamChunk struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func HandlePromptStream(c *websocket.Conn) {
	userIDRaw := c.Locals("user_id")
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		log.Println("❌ Invalid or missing user_id in WebSocket")
		c.WriteMessage(websocket.TextMessage, []byte("⚠️ Not authenticated"))
		return
	}

	log.Println("✅ WebSocket connected for user ID:", userID)

	defer func() {
		log.Println("Closing WebSocket")
		c.Close()
	}()

	db := database.Get()

	// 🔹 1. Load or create user stats
	var userStats models.UserStats
	if err := db.Where("user_id = ?", userID).First(&userStats).Error; err != nil {
		userStats = models.UserStats{
			UserID:        userID,
			QuestionsUsed: 0,
		}
		if err := db.Create(&userStats).Error; err != nil {
			log.Println("❌ Failed to create user stats:", err)
			c.WriteMessage(websocket.TextMessage, []byte("⚠️ Unable to track usage"))
			return
		}
	}

	// 🔹 2. Get JSON message from frontend
	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("WebSocket read error:", err)
		return
	}

	var req GenerateRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		log.Println("Invalid JSON input:", err)
		c.WriteMessage(websocket.TextMessage, []byte("⚠️ Invalid input"))
		return
	}

	// 🔹 4. Ensure streaming is enabled
	req.Stream = true
	req.Stop = []string{"</s>"} // optional, for clean stop

	// 🔹 5. Call Ollama API
	body, err := json.Marshal(req)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://ollama:11434"
	}

	log.Println("🔧 OLLAMA_URL from env:", os.Getenv("OLLAMA_URL"))

	resp, err := http.Post(ollamaURL+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Println("Ollama request error:", err)
		c.WriteMessage(websocket.TextMessage, []byte("⚠️ Failed to reach model"))
		return
	}
	defer resp.Body.Close()

	// 🔹 6. Stream response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		var chunk OllamaStreamChunk
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			log.Println("Stream decode error:", err)
			continue
		}

		// Send each chunk over WebSocket
		if err := c.WriteMessage(websocket.TextMessage, []byte(chunk.Response)); err != nil {
			log.Println("WebSocket write error:", err)
			break
		}

		if chunk.Done {
			break
		}
	}

	// 🔹 7. Update usage
	userStats.QuestionsUsed++
	if err := db.Save(&userStats).Error; err != nil {
		log.Println("❌ Failed to update user stats:", err)
	}

	// 🔹 8. Save prompt to DB
	hashedPrompt := hashPrompt(req.Prompt)

	if err := db.Create(&models.Question{
		UserID: userID,
		Model:  req.Model,
		Prompt: hashedPrompt,
	}).Error; err != nil {
		log.Println("❌ Failed to save question:", err)
	}
}

func AskPageHandler(c *fiber.Ctx) error {
	userID := authentication.GetUserID(c)

	// Get user info
	var user models.User
	if err := database.Get().First(&user, userID).Error; err != nil {
		return c.Redirect("/")
	}

	// Get or create user stats
	var userStats models.UserStats
	if err := database.Get().Where("user_id = ?", userID).First(&userStats).Error; err != nil {
		userStats = models.UserStats{
			UserID:        userID,
			QuestionsUsed: 0,
			DailyResetAt:  time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		database.Get().Create(&userStats)
	}

	// Reset daily quota if needed
	now := time.Now()
	if now.After(userStats.DailyResetAt) {
		userStats.QuestionsUsed = 0
		userStats.DailyResetAt = now.AddDate(0, 0, 1).Truncate(24 * time.Hour)
		database.Get().Save(&userStats)
	}

	const questionsLimit = 25 // Community version limit per day
	questionsRemaining := questionsLimit - userStats.QuestionsUsed
	if questionsRemaining < 0 {
		questionsRemaining = 0
	}

	// Load available models for dropdown
	modelGroups := services.GetAvailableModels()

	return c.Render("ask", fiber.Map{
		"Title":              "Chat",
		"User":               user,
		"QuestionsRemaining": questionsRemaining,
		"QuestionsLimit":     questionsLimit,
		"HasQuota":           questionsRemaining > 0,
		"ModelGroups":        modelGroups,
	}, "layouts/main")
}

func AskHandler(c *fiber.Ctx) error {
	userID := authentication.GetUserID(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).SendString("Not authenticated")
	}

	prompt := c.FormValue("prompt")
	model := c.FormValue("model")
	if model == "" {
		model = "phi"
	}
	if prompt == "" {
		return c.Status(400).SendString("Prompt cannot be empty")
	}

	// Create request to Ollama
	req := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
		Stop:   []string{"</s>"},
	}

	payload, _ := json.Marshal(req)
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://ollama:11434"
	}

	log.Println("🔧 OLLAMA_URL from env:", os.Getenv("OLLAMA_URL"))

	resp, err := http.Post(ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(payload))

	if err != nil {
		return c.Status(500).SendString("Error connecting to Ollama")
	}

	// Save prompt (trimmed for privacy)
	promptToSave := prompt
	if len(promptToSave) > 100 {
		promptToSave = promptToSave[:100]
	}
	_ = database.Get().Create(&models.Question{
		UserID: userID,
		Model:  model,
		Prompt: promptToSave,
	})

	// Stream response
	responseHTML := `
	<div id="answer-container" class="mt-4 bg-gray-100 rounded-lg p-4 text-gray-800">
		<div id="ai-response"></div>
	</div>
	<script>
		document.querySelector('textarea[name="prompt"]').value = "";
	</script>
	`

	c.Set("Content-Type", "text/html")
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		w.WriteString(responseHTML)
		w.Flush()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var chunk map[string]interface{}
			if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
				continue
			}
			if text, ok := chunk["response"].(string); ok {
				w.WriteString("<script>document.getElementById('ai-response').innerHTML += " +
					strconv.Quote(text) + ";</script>")
				w.Flush()
			}
		}
		resp.Body.Close()
	})

	return nil
}

func hashPrompt(prompt string) string {
	h := sha256.New()
	h.Write([]byte(prompt))
	return hex.EncodeToString(h.Sum(nil))
}
