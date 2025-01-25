package routes

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/choccops/auth/internal/models"
)

func Signin(c *fiber.Ctx) error {
	request := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&request); err != nil {
		fmt.Fprintf(os.Stderr, "BodyParser: %v\n", err)
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	if request.Username == "" {
		return c.JSON(fiber.Map{
			"oops": "missing username",
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

	postgres, ok := c.Locals("postgres").(*pgxpool.Pool)
	if !ok {
		fmt.Fprintf(os.Stderr, "Couldn't find postgres\n")
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	query := `SELECT * FROM users WHERE LOWER(username) = LOWER($1) LIMIT 1`
	row, err := postgres.Query(context.Background(), query, request.Username)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Query: %v\n", err)
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.User])

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(fiber.Map{
				"oops": "username not found",
			})
		}

		fmt.Fprintf(os.Stderr, "CollectOneRow: %v\n", err)
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if err != nil {
		fmt.Fprintf(os.Stderr, "CompareHashAndPassword: %v\n", err)
		return c.JSON(fiber.Map{
			"oops": "wrong password",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	privateKey, ok := c.Locals("private-key").(*rsa.PrivateKey)
	if !ok {
		fmt.Fprintf(os.Stderr, "Couldn't find privateKey")
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		fmt.Fprintf(os.Stderr, "SignedString: %v\n", err)
		return c.JSON(fiber.Map{
			"oops": "something went wrong",
		})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
	})

}
