package controller

import (
	"book/model"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func GetAll(db *gorm.DB, c *fiber.Ctx, jwtSecret []byte) error {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, no token provided",
		})
	}
	// check JWT and Decode
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		// check Signing Method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	// check truth JWT
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid token",
		})
	}

	// get Claims an check user_id
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid token claims",
		})
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid user_id format",
		})
	}

	var books []model.Book
	if err := db.Where("user_id = ?", int(userID)).Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch books for the user",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Books fetched successfully",
		"total_books": len(books),
		"books":       books,
	})
}
