package controller

import (
	"book/model"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(db *gorm.DB, c *fiber.Ctx, jwtSecret []byte) error {
	// Parse JSON request body
	req := new(model.User)
	if err := c.BodyParser(req); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	// Query user from the database
	var user model.User
	result := db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		log.Println("User not found or query error:", result.Error)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "email or password are wrong"})
	}

	// Compare hash password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Println("Password comparison failed:", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "password is incorrect"})
	}

	// Create JWT token
	data := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 3).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, data)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Failed to sign JWT token:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create jwt token"})
	}

	// Set JWT token in HTTP cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.Expires = time.Now().Add(time.Hour * 3)
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.SameSite = "Lax"
	c.Cookie(cookie)

	log.Println("Login successful for user ID:", user.ID)
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Login successful"})
}
