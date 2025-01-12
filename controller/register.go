package controller

import (
	"book/model"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"github.com/joho/godotenv"
	"os"
	"log"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gopkg.in/gomail.v2"
)

func Register(db *gorm.DB, c *fiber.Ctx) error {
	req := new(model.User)
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	var existingUser model.User
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "email already registered"})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "cannot hash password"})
	}
	req.Password = string(hashPassword)

	// verification Token
	token := generateVerificationToken()
	req.VerificationToken = token
	req.IsEmailVerified = false


	if err := db.Create(req).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "cannot register new user"})
	}

	if err := sendVerificationEmail(req.Email, token); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to send verification email"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message":"Check Verify in your Email."})
}

// create token
func generateVerificationToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// send authen email
func sendVerificationEmail(email, token string) error {
	// .env
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    gmail := os.Getenv("GMAIL")
    password := os.Getenv("GPASS")
	host := os.Getenv("DOMAIN")

    if gmail == "" || password == "" {
        log.Fatal("Environment variables GMAIL or GPASS are missing")
    }

    m := gomail.NewMessage()
    m.SetHeader("From", gmail)
    m.SetHeader("To", email)
    m.SetHeader("Subject", "Email Verification")
    m.SetBody("text/plain", fmt.Sprintf("Your verification token: %s/verify?token=%s", host, token))

    // ตั้งค่า SMTP สำหรับ Gmail
    d := gomail.NewDialer("smtp.gmail.com", 587, gmail, password)

    if err := d.DialAndSend(m); err != nil {
        log.Println("Failed to send email:", err)
        return err
    }

    log.Println("Verification email sent successfully to:", email)
    return nil
}
