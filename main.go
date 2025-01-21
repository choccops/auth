package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	env := os.Getenv("ENV")
	if "" == env {
		godotenv.Load()
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		message := os.Getenv("ASDOMARE")
		return c.JSON(fiber.Map{
			"message": message,
		})
	})

	app.Listen(":3000")
}
