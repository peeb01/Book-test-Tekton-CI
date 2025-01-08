package controller

import (
	"book/model"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewBook(db *gorm.DB, c *fiber.Ctx) error {
	book := new(model.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "some data is missing"})
	}

	var exbook model.Book
	if err := db.Where("isbn = ?, title = ?, author = ?", book.ISBN, book.Title, book.Author).First(&exbook).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "book are in library"})
	} else if !errors.Is(err, gorm.ErrRecordNotFound){
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error checking existing book"})
	}

	if err := db.Create(book).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create new book to the database"})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "successful"})
}