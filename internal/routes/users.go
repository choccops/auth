package routes

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/choccops/auth/internal/models"
)

func Users(c *fiber.Ctx) error {
	request := struct {
		Token string `json:"token"`
	}{}

	if err := c.BodyParser(&request); err != nil {
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	publicKey, ok := c.Locals("public-key").(*rsa.PublicKey)
	if !ok {
		fmt.Fprintf(os.Stderr, "Couldn't find privateKey")
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse: %v\n", err)
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	if !token.Valid {
		return c.JSON(fiber.Map{
			"oops": "invalid token",
		})
	}

	postgres, ok := c.Locals("postgres").(*pgxpool.Pool)
	if !ok {
		fmt.Fprintf(os.Stderr, "Couldn't find postgres\n")
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	rows, err := postgres.Query(context.Background(), "select * from users")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Query: %v\n", err)
		return c.JSON(fiber.Map{
			"oops": err,
		})
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])

	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows: %v\n", err)
		return c.JSON(fiber.Map{
			"oops": err,
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}
