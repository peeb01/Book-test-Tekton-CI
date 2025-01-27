package controller

import (
	"book/model"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func SummaryBooks(db *gorm.DB, c *fiber.Ctx, jwtSecret []byte) error {

	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized, no token provided"})
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized, invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized, invalid token claims"})
	}

	userID := int(claims["user_id"].(float64))

	bookTypes := []string{"Light Novel", "Manga", "Knowledge"}

	summary := make(map[string]int)
	for _, bookType := range bookTypes {
		summary[bookType] = 0
	}
	summary["Other"] = 0 

	var books []model.Book
	if err := db.Where("user_id = ?", userID).Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch books"})
	}

	totalBooks := len(books)

	for _, book := range books {
		if _, exists := summary[book.BookType]; exists {
			summary[book.BookType]++
		} else {
			summary["Other"]++
		}
	}

	return c.JSON(fiber.Map{
		"total_books": totalBooks,
		"summary":     summary,
	})
}
