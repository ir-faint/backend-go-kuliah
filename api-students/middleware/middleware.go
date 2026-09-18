package middleware

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"api-students/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var payloadMethods = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

// RequireJSON memastikan header Content-Type bernilai application/json untuk permintaan bernilai body.
func RequireJSON(c *fiber.Ctx) error {
	if payloadMethods[c.Method()] {
		if !strings.HasPrefix(c.Get("Content-Type"), fiber.MIMEApplicationJSON) {
			return helper.Fail(c, fiber.StatusUnsupportedMediaType, "Content-Type harus application/json")
		}
	}
	return c.Next()
}

// RequestLogger mencatat setiap request ke log terstruktur.
func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next() // serahkan ke middleware/handler berikutnya

		requestID, _ := c.Locals("requestid").(string)

		logger.Info("http_request",
			slog.String("request_id", requestID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", c.IP()),
		)

		return err
	}
}

// Register memasang seluruh middleware yang berlaku untuk semua route.
func Register(app *fiber.App, logger *slog.Logger) {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")
	}
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	app.Use(requestid.New()) // 1. beri setiap request satu ID unik
	app.Use(recover.New())   // 2. tangkap panic agar server tidak mati
	app.Use(helmet.New())    // 3. pasang header keamanan dasar
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	})) // 4. atur Cross-Origin Resource Sharing
	app.Use(RequestLogger(logger)) // 5. catat setiap request
}
