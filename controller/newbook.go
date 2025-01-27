package controller

import (
	"book/model"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func NewBook(db *gorm.DB, c *fiber.Ctx, jwtSecret []byte) error {
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

	book := new(model.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "some data is missing or invalid"})
	}

	var existingBook model.Book
	if err := db.Where("isbn = ? AND user_id = ?", book.ISBN, userID).First(&existingBook).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "book with this ISBN already exists for the user"})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error checking existing book"})
	}

	book.UserID = uint(userID)

	if err := db.Create(book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create new book in the database"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "book added successfully"})
}
