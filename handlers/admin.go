package handlers

import (
	"LibreAI/internal/database"
	"LibreAI/models"

	"github.com/gofiber/fiber/v2"
)

func AdminModelPage(c *fiber.Ctx) error {
	var models []models.Model
	database.Get().Order("category").Find(&models)

	return c.Render("admin/models", fiber.Map{
		"Title":  "Model Manager",
		"Models": models,
	}, "layouts/main")
}

func AdminModelPartial(c *fiber.Ctx) error {
	var models []models.Model
	database.Get().Order("category").Find(&models)

	return c.Render("admin/_model_list", fiber.Map{
		"Models": models,
	})
}

func AdminModelCreate(c *fiber.Ctx) error {
	model := new(models.Model)
	if err := c.BodyParser(model); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	model.IsActive = true
	if err := database.Get().Create(model).Error; err != nil {
		return c.Status(500).SendString("Failed to create model")
	}

	return AdminModelPartial(c)
}

func AdminModelToggle(c *fiber.Ctx) error {
	id := c.Params("id")
	var model models.Model
	if err := database.Get().First(&model, id).Error; err != nil {
		return c.Status(404).SendString("Model not found")
	}

	model.IsActive = !model.IsActive
	database.Get().Save(&model)

	return AdminModelPartial(c)
}

func AdminModelDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.Get().Delete(&models.Model{}, id).Error; err != nil {
		return c.Status(500).SendString("Failed to delete model")
	}

	return AdminModelPartial(c)
}
