package routes

import (
	"book/controller"
	"fmt"

	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func Routes(db *gorm.DB){
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	JWT := []byte(os.Getenv("JWT_SECRET"))
	fmt.Println(JWT)

	app := fiber.New()
	app.Use(fiberlog.New(fiberlog.Config{
		Format:     "[${time}] ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))
	
	app.Post("/login", func(c *fiber.Ctx) error{
		return controller.Login(db, c, JWT)
	})
	app.Post("/register", func(c *fiber.Ctx) error{
		return controller.Register(db, c)
	})
	app.Post("/newbook", func(c *fiber.Ctx) error{
		return controller.NewBook(db, c)
	})
	app.Get("/verify", func(c *fiber.Ctx) error {
		return controller.VerifyEmail(db, c)
	})
	app.Get("/summary", func(c *fiber.Ctx) error {
		return controller.SummaryBooks(db, c)
	})
	app.Get("/search", func(c *fiber.Ctx) error {
		return controller.SearchBook(db, c)
	})
	app.Get("/book", func(c *fiber.Ctx) error {
		return controller.GetAll(db, c)
	})

	app.Listen(":8000")

}
