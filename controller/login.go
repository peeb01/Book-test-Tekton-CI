package controller

import (
	"book/model"
	// "fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(db *gorm.DB, c *fiber.Ctx, jwtSecret []byte) error {
	req := new(model.User)
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	var user model.User
	result := db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "email or password are wrong"})
	}

	// compare hash password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "password wrong"})
	}

	// create jwt
	data := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 3).Unix(),
	}
	// fmt.Println("JWT Secret as string:", string(jwtSecret))
	// fmt.Printf("JWT Secret as bytes: %v\n", jwtSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, data) 
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create jwt token"})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.Expires = time.Now().Add(time.Hour * 3)
	cookie.HTTPOnly = true
	cookie.Secure = true 
	cookie.SameSite = "Lax" 
	c.Cookie(cookie)
	
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Login successful"})
}