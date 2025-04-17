package services

import (
	"LibreAI/internal/database"
	"LibreAI/models"
)

type ModelGroup struct {
	Label   string
	Options []models.Model
}

// GetAvailableModelsByTier returns models grouped by category (e.g., small, medium) for the given tier
func GetAvailableModels() []ModelGroup {
	db := database.Get()
	var all []models.Model
	db.Where("is_active = ?", true).Order("category, name").Find(&all)

	groups := map[string][]models.Model{}
	for _, model := range all {
		groups[model.Category] = append(groups[model.Category], model)
	}

	// Convert to slice of groups with labels
	ordered := []ModelGroup{}

	if small := groups["small"]; len(small) > 0 {
		ordered = append(ordered, ModelGroup{
			Label:   "🟢 Small Models (Fastest Responses)",
			Options: small,
		})
	}
	if medium := groups["medium"]; len(medium) > 0 {
		ordered = append(ordered, ModelGroup{
			Label:   "🟡 Medium Models",
			Options: medium,
		})
	}
	if large := groups["large"]; len(large) > 0 {
		ordered = append(ordered, ModelGroup{
			Label:   "🔴 Large Models",
			Options: large,
		})
	}

	return ordered
}
