package services

import (
	"LibreAI/internal/database"
	"LibreAI/models"
	"errors"
)

var (
	ErrQuotaExceeded = errors.New("daily question quota exceeded")
)

func GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := database.Get().First(&user, userID).Error
	return &user, err
}
