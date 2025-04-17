// models/user_stats.go
package models

import "time"

type UserStats struct {
	ID            int `gorm:"primaryKey"`
	UserID        int `gorm:"unique"`
	QuestionsUsed int `gorm:"default:0"`
	DailyResetAt  time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
