package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type User struct {
	UserID    int       `json:"userID"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	/** local development */
	if "" == os.Getenv("ENV") {
		godotenv.Load()
	}

	/** postgres connection */
	postgres, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URI"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer postgres.Close()

	/** fiber */
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		message := os.Getenv("ASDOMARE")
		return c.JSON(fiber.Map{
			"message": message,
		})
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		rows, _ := postgres.Query(context.Background(), "select * from users")
		users, err := pgx.CollectRows(rows, pgx.RowTo[User])
		if err != nil {
			return c.JSON(fiber.Map{
				"oops": err,
			})
		}

		return c.JSON(fiber.Map{
			"users": users,
		})
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		request := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		if err := c.BodyParser(&request); err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"username": request.Username,
			"password": request.Password,
		})
	})

	app.Listen(":3000")
}
