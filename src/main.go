package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProxyLog struct {
	gorm.Model
	Request  string `json:"request"`
	Response string `json:"response"`
}

var db *gorm.DB

func main() {
	godotenv.Load()
	app := fiber.New()

	dsn := fmt.Sprintf("postgresql://cosmin:%s@dapper-ape-12047.8nj.cockroachlabs.cloud:26257/%s?sslmode=verify-full", os.Getenv("SQL_USER_PASSWORD"), os.Getenv("DATABASE_NAME"))

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	err = db.AutoMigrate(&ProxyLog{})
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
			return c.SendFile("./ratelimit.html")
		},
	}))

	app.All("/*", reverseProxyHandler)

	log.Fatal(app.Listen(":8080"))
}
