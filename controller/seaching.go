package controller

import (
	"book/model"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SearchBook(db *gorm.DB, c *fiber.Ctx) error {
    query := c.Query("q")
    if query == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "search query is missing"})
    }

    processedQuery := normalize(query)

    var books []model.Book
    err := db.Where(`
        replace(lower(Title), ' ', '') like ? or
        replace(lower(Author), ' ', '') like ? or
        replace(lower(ISBN), ' ', '') like ?`,
        "%"+processedQuery+"%", "%"+processedQuery+"%", "%"+processedQuery+"%",
    ).Find(&books).Error

    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error fetch"})
    }

    return c.Status(http.StatusOK).JSON(books)
}

func normalize(input string) string {
	return strings.ToLower(strings.ReplaceAll(input, " ", ""))
}