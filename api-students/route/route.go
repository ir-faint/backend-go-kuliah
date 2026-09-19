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

type Dependencies struct {
	Pool           *pgxpool.Pool
	StudentService *service.StudentService
	AuthService    *service.AuthService
	JWT            *helper.JWTManager
}

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
func RegisterRoutes(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	// Public endpoints
	api.Get("/health", HealthCheck(deps.Pool))

	// Auth endpoints
	auth := api.Group("/auth", middleware.RequireJSON)
	auth.Post("/register", deps.AuthService.Register)
	auth.Post("/login", middleware.LoginRateLimiter(), deps.AuthService.Login)
	auth.Post("/refresh", deps.AuthService.Refresh)
	auth.Post("/logout", deps.AuthService.Logout)
	auth.Get("/me", middleware.RequireAuth(deps.JWT), deps.AuthService.Me)

	// Protected endpoints (All students endpoints require auth)
	students := api.Group("/students", middleware.RequireAuth(deps.JWT), middleware.RequireJSON)
	students.Get("/", deps.StudentService.ListStudents)
	students.Get("/:id", deps.StudentService.GetStudent)
	students.Post("/", deps.StudentService.CreateStudent)
	students.Put("/:id", deps.StudentService.ReplaceStudent)
	students.Patch("/:id", deps.StudentService.PatchStudent)
	students.Delete("/:id", deps.StudentService.DeleteStudent)
}
