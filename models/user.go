package models

import "time"

type User struct {
	ID         int    `gorm:"primaryKey"`
	Username   string `gorm:"unique"`
	Image      string
	Provider   string
	ProviderID string `gorm:"column:provider_id;unique"`
	Name       string
	Email      string
	AvatarURL  string
	Role       string `gorm:"default:user"`
	CreatedAt  time.Time
}
