package controller

import (
	"book/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SummaryBooks(db *gorm.DB, c *fiber.Ctx) error {
	bookTypes := []string{"Light Novel", "Manga", "Knowledge"}
	
	summary := make(map[string]int)
	for _, bookType := range bookTypes {
		summary[bookType] = 0
	}
	var books []model.Book
	if err := db.Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch books"})
	}

	totalBooks := len(books)
	// counting book
	for _, book := range books {
		if _, exists := summary[book.BookType]; exists {
			summary[book.BookType]++
		} else {
			// if not in booktype
			summary["Other"]++
		}
	}

	return c.JSON(fiber.Map{
		"total_books": totalBooks,
		"summary":     summary,
	})
}
