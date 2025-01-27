package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Hello(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Hello World"})
}