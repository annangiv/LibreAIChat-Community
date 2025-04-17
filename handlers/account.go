// handlers/account.go
package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"LibreAI/authentication"
	"LibreAI/internal/database"
	"LibreAI/models"
	"LibreAI/utils"
)

// AccountHandler renders the account page with user info and stats
func AccountHandler(c *fiber.Ctx) error {
	userID := authentication.GetUserID(c)
	if userID == 0 {
		return c.Redirect("/")
	}

	db := database.Get()

	// Get user
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Redirect("/")
	}

	// Get or create user stats
	var userStats models.UserStats
	if err := db.Where("user_id = ?", userID).First(&userStats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userStats = models.UserStats{
				UserID:        userID,
				QuestionsUsed: 0,
				DailyResetAt:  time.Now().Add(24 * time.Hour),
			}
			if err := db.Create(&userStats).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to create user stats")
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to load user stats")
		}
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	var (
		questionsToday     int64
		questionsThisMonth int64
		questionsTotal     int64
	)

	db.Model(&models.Question{}).Where("user_id = ? AND created_at >= ?", userID, startOfDay).Count(&questionsToday)
	db.Model(&models.Question{}).Where("user_id = ? AND created_at >= ?", userID, startOfMonth).Count(&questionsThisMonth)
	db.Model(&models.Question{}).Where("user_id = ?", userID).Count(&questionsTotal)

	// Model usage stats
	type ModelUsageResult struct {
		ModelName string
		Count     int
	}
	var modelUsageResults []ModelUsageResult
	db.Raw(`
		SELECT model AS model_name, COUNT(*) AS count 
		FROM questions 
		WHERE user_id = ? 
		GROUP BY model 
		ORDER BY count DESC 
		LIMIT 5
	`, userID).Scan(&modelUsageResults)

	modelUsage := make([]map[string]interface{}, 0, len(modelUsageResults))
	for _, r := range modelUsageResults {
		modelUsage = append(modelUsage, map[string]interface{}{
			"ModelName": r.ModelName,
			"Count":     r.Count,
		})
	}

	return c.Render("account", fiber.Map{
		"Title": "My Account",
		"User": fiber.Map{
			"ID":        user.ID,
			"Name":      user.Name,
			"Email":     utils.MaskEmail(user.Email),
			"AvatarURL": user.AvatarURL,
		},
		"Stats": fiber.Map{
			"QuestionsToday":     questionsToday,
			"QuestionsThisMonth": questionsThisMonth,
			"QuestionsTotal":     questionsTotal,
			"ModelUsage":         modelUsage,
		},
	}, "layouts/main")
}

func DeleteAccount(c *fiber.Ctx) error {
	userID := authentication.GetUserID(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// ⚠️ Get user first BEFORE deleting
	db := database.Get()
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
	}

	if err := deleteUserAndData(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	authentication.AuthLogout(c)
	return c.Redirect("/goodbye")
}

func deleteUserAndData(userID int) error {
	db := database.Get()
	tx := db.Begin()

	// Delete related tables
	tables := []string{"user_stats", "questions"}
	for _, table := range tables {
		if err := tx.Table(table).Where("user_id = ?", userID).Delete(nil).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete from %s: %w", table, err)
		}
	}

	// Delete user
	if err := tx.Delete(&models.User{}, userID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return tx.Commit().Error
}

func MeHandler(c *fiber.Ctx) error {
	userID := authentication.GetUserID(c)
	if userID == 0 {
		// Return the default login buttons
		return c.Render("components/login_buttons", fiber.Map{}, "")
	}

	// Get user info
	var user models.User
	if err := database.Get().First(&user, userID).Error; err != nil {
		return c.Render("components/login_buttons", fiber.Map{}, "")
	}

	// Return the logged-in menu
	return c.Render("components/user_menu", fiber.Map{
		"User": user,
	}, "")
}
