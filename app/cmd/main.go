package main

import (
	"os"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ready bool = false

var requestCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total Number of HTTP requests",
	},
)

func init() {
	prometheus.MustRegister(requestCounter)
}

func main() {
	app := fiber.New()

	// simulate startup delay
	go func() {
		time.Sleep(15 * time.Second)
		ready = true
	}()

	app.Get("/", func(c *fiber.Ctx) error {

		requestCounter.Inc()

		return c.JSON(fiber.Map{
			"message":      "Go API Running",
			"environment":  os.Getenv("APP_ENV"),
			"feature_flag": os.Getenv("FEATURE_FLAG"),
			"db_user":      os.Getenv("DATABASE_USER"),
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

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
