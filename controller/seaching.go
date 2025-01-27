package controller

import (
	"book/model"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func SearchBook(db *gorm.DB, c *fiber.Ctx, jwtSecret []byte) error {
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

	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "search query is missing"})
	}

	processedQuery := normalize(query)

	var books []model.Book
	err = db.Where(`
		(user_id = ?) AND (
			replace(lower(Title), ' ', '') like ? OR
			replace(lower(Author), ' ', '') like ? OR
			replace(lower(ISBN), ' ', '') like ?
		)`,
		userID, "%"+processedQuery+"%", "%"+processedQuery+"%", "%"+processedQuery+"%",
	).Find(&books).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching books"})
	}

	return c.Status(fiber.StatusOK).JSON(books)
}

// normalize function to manage string in Query
func normalize(input string) string {
	return strings.ToLower(strings.ReplaceAll(input, " ", ""))
}
