package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"github.com/choccops/auth/internal/routes"
)

const Service = "auth"
const Version = "0.1.16"

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

	/** migrations */
	db, err := sql.Open("pgx", os.Getenv("POSTGRES_URI"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error wrapping pgx: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting dialeact: %v\n", err)
		os.Exit(1)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		fmt.Fprintf(os.Stderr, "Error during migrations: %v\n", err)
		os.Exit(1)
	}

	/** fiber */
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("postgres", postgres)
		c.Locals("private-key", privateKey)
		c.Locals("public-key", publicKey)
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": Service,
			"version": Version,
		})
	})

	app.Get("/users", routes.Users)
	app.Post("/signup", routes.Signup)
	app.Get("/signin", routes.Signin)

	app.Listen(":3000")
}
