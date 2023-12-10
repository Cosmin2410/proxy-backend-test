package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Cosmin2410/proxy-backend-test/src/handler"
	"github.com/Cosmin2410/proxy-backend-test/src/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	dsn := fmt.Sprintf(
		"postgresql://cosmin:%s@dapper-ape-12047.8nj.cockroachlabs.cloud:26257/%s?sslmode=verify-full",
		os.Getenv("SQL_USER_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)

	if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Fatal("Failed to connect database", err)
	}

	err = db.AutoMigrate(&model.SaveLog{})
	if err != nil {
		log.Fatal("Auto migrate failed", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Tibi, reverse-proxy is up!")
	})

	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 60 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendFile("../ratelimit.html")
		},
	}))

	h := &handler.DBCreate{DB: db}
	app.All("/*", h.ReverseProxyHandler)

	log.Fatal(app.Listen(":8080"))
}
