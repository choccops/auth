package routes

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {
	request := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	if request.Username == "" {
		return c.JSON(fiber.Map{
			"oops": "missing username",
		})
	}

	postgres, ok := c.Locals("postgres").(*pgxpool.Pool)
	if !ok {
		fmt.Fprintf(os.Stderr, "Couldn't find postgres\n")
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`

	var exists bool
	err := postgres.QueryRow(context.Background(), query, request.Username).Scan(&exists)

	if err != nil {
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	if exists {
		return c.JSON(fiber.Map{
			"oops": "username already exists",
		})
	}

	match, _ := regexp.MatchString("^[a-zA-Z._]{4,16}$", request.Username)
	if match == false {
		return c.JSON(fiber.Map{
			"oops": "username is invalid",
		})
	}

	if request.Password == "" {
		return c.JSON(fiber.Map{
			"oops": "missing password",
		})
	}

	match, _ = regexp.MatchString(`^.{8,64}$`, request.Password)

	if match == false {
		return c.JSON(fiber.Map{
			"oops": "password is invalid",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	query = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`

	var id int
	err = postgres.QueryRow(context.Background(), query, request.Username, hash).Scan(&id)

	if err != nil {
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	return c.JSON(fiber.Map{
		"id": id,
	})
}
