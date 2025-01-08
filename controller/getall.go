package controller

import (
	"book/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAll(db *gorm.DB, c *fiber.Ctx) error {
	var books []model.Book
	if err := db.Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch books",
		})
	}
	return c.JSON(books)
}