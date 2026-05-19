package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

var ready bool = false

func main() {
	app := fiber.New()

	// simulate startup delay
	go func() {
		time.Sleep(15 * time.Second)
		ready = true
	}()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Go API Running",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Get("/ready", func(c *fiber.Ctx) error {
		if !ready {
			return c.Status(503).SendString("Not Ready")
		}

		return c.SendString("Ready")
	})

	app.Get("/crash", func(c *fiber.Ctx) error {
		os.Exit(1)
		return nil
	})

	app.Listen(":5001")
}
