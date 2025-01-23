package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

const Service = "auth"
const Version = "0.1.1"

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func main() {
	/** local development */
	if "" == os.Getenv("ENV") {
		godotenv.Load()
	}

	/** load keys */
	privateKeyBytes, err := os.ReadFile(os.Getenv("PRIVATE_KEY"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open private key: %v\n", err)
		os.Exit(1)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse private key: %v\n", err)
		os.Exit(1)
	}

	publicKeyBytes, err := os.ReadFile(os.Getenv("PUBLIC_KEY"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open public key: %v\n", err)
		os.Exit(1)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

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
		return c.JSON(fiber.Map{
			"service": Service,
			"version": Version,
		})
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		request := struct {
			Token string `json:"token"`
		}{}

		if err := c.BodyParser(&request); err != nil {
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

		rows, err := postgres.Query(context.Background(), "select * from users")

		if err != nil {
			fmt.Fprintf(os.Stderr, "Query: %v\n", err)
			return c.JSON(fiber.Map{
				"oops": err,
			})
		}

		users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])

		if err != nil {
			fmt.Fprintf(os.Stderr, "CollectRows: %v\n", err)
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
			return c.JSON(fiber.Map{
				"oops": "something went wrong",
			})
		}

		if request.Username == "" {
			return c.JSON(fiber.Map{
				"oops": "missing username",
			})
		}

		query := `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`

		var exists bool
		err = postgres.QueryRow(context.Background(), query, request.Username).Scan(&exists)

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
	})

	app.Get("/signin", func(c *fiber.Ctx) error {
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

		query := `SELECT * FROM users WHERE LOWER(username) = LOWER($1) LIMIT 1`
		row, err := postgres.Query(context.Background(), query, request.Username)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Query: %v\n", err)
			return c.JSON(fiber.Map{
				"oops": "something went wrong",
			})
		}

		user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[User])

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

	})

	app.Listen(":3000")
}
