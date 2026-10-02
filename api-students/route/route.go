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
	Permissions    *helper.PermissionSet
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

	// Protected endpoints (Students domain)
	students := api.Group("/students", middleware.RequireAuth(deps.JWT), middleware.RequireJSON)
	students.Get("/", middleware.RequirePermission(deps.Permissions, "student:list"), deps.StudentService.ListStudents)
	students.Post("/", middleware.RequirePermission(deps.Permissions, "student:create"), deps.StudentService.CreateStudent)
	students.Delete("/:id", middleware.RequirePermission(deps.Permissions, "student:delete"), deps.StudentService.DeleteStudent)

	// Endpoints guarded by service-level ownership checks (CanAccessStudent)
	students.Get("/:id", deps.StudentService.GetStudent)
	students.Put("/:id", deps.StudentService.ReplaceStudent)
	students.Patch("/:id", deps.StudentService.PatchStudent)

	// Protected endpoints (User Management domain)
	users := api.Group("/users", middleware.RequireAuth(deps.JWT), middleware.RequireJSON)
	users.Put("/:id/role", middleware.RequirePermission(deps.Permissions, "role:assign"), deps.AuthService.AssignRole)
	users.Delete("/:id", middleware.RequirePermission(deps.Permissions, "user:delete"), deps.AuthService.DeleteUser)
}
