package main

import (
	"context"
	"log"
	"strings"
	"time"

	"api-students/config"
	"api-students/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var payloadMethods = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

func requireJSON(c *fiber.Ctx) error {
	if payloadMethods[c.Method()] {
		if !strings.HasPrefix(c.Get("Content-Type"), fiber.MIMEApplicationJSON) {
			return fail(c, fiber.StatusUnsupportedMediaType, "Content-Type harus application/json")
		}
	}
	return c.Next()
}

func main() {
	config.LoadEnv()

	ctx := context.Background()
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Printf("Peringatan koneksi database: %v", err)
	} else {
		defer pool.Close()
	}

	app := fiber.New(fiber.Config{
		AppName: "API Students - Praktikum Backend",
	})

	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		if pool == nil {
			return fail(c, fiber.StatusServiceUnavailable, "database tidak terhubung")
		}

		pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := pool.Ping(pingCtx); err != nil {
			return fail(c, fiber.StatusServiceUnavailable, "database tidak tersedia")
		}

		return ok(c, "server berjalan", fiber.Map{"timestamp": time.Now()})
	})
	s := api.Group("/students", requireJSON)

	s.Get("/", listStudents)
	s.Get("/:id", getStudent)
	s.Post("/", createStudent)
	s.Put("/:id", replaceStudent)
	s.Patch("/:id", patchStudent)
	s.Delete("/:id", deleteStudent)

	app.Use(func(c *fiber.Ctx) error {
		return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	log.Fatal(app.Listen(":3000"))
}

