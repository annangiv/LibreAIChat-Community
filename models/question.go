// models/question.go
package models

import "time"

type Question struct {
	ID        int    `gorm:"primaryKey"`
	UserID    int    `gorm:"index"`
	Model     string `gorm:"size:50"` // Model used (e.g., "phi", "gemma", etc.)
	Prompt    string `gorm:"size:64"` // SHA-256 hash (hex is 64 chars)
	CreatedAt time.Time
}
