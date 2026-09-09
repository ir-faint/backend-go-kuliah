package route

import (
	"context"
	"time"

	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthCheck mengembalikan handler fiber untuk memeriksa kesehatan database.
func HealthCheck(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if pool == nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable, "database tidak terhubung")
		}

		pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := pool.Ping(pingCtx); err != nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable, "database tidak tersedia")
		}

		return helper.Success(c, "server berjalan", fiber.Map{"timestamp": time.Now()})
	}
}

// RegisterRoutes mendaftarkan seluruh endpoint API ke instansi Fiber.
func RegisterRoutes(app *fiber.App, studentService *service.StudentService, pool *pgxpool.Pool) {
	api := app.Group("/api/v1")
	api.Get("/health", HealthCheck(pool))

	students := api.Group("/students", middleware.RequireJSON)

	students.Get("/", studentService.ListStudents)
	students.Get("/:id", studentService.GetStudent)
	students.Post("/", studentService.CreateStudent)
	students.Put("/:id", studentService.ReplaceStudent)
	students.Patch("/:id", studentService.PatchStudent)
	students.Delete("/:id", studentService.DeleteStudent)
}
