package controller

import (
	"book/model"
	// "net/http"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func VerifyEmail(db *gorm.DB, c *fiber.Ctx) error {
	token := c.Query("token")

	var user model.User
	if err := db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	user.IsEmailVerified = true
	user.VerificationToken = ""
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to verify email"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "email verified successfully"})
}
