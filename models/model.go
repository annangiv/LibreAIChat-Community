package models

import (
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name         string `json:"name"`
	Identifier   string `json:"key" gorm:"unique"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	RequiredTier string `json:"required_tier"`
	IsActive     bool   `json:"is_active"`
}
