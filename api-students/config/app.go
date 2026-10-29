package config

import (
	"errors"
	"log/slog"

	"api-students/helper"
	"api-students/middleware"
	"api-students/app/model"
	"api-students/route"

	"github.com/gofiber/fiber/v2"
)

// NewApp mendaftarkan route. File ini adalah tempat seluruh bagian bertemu.
func NewApp(logger *slog.Logger, deps route.Dependencies) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      GetEnv("APP_NAME", "API Students - Praktikum Backend"),
		BodyLimit:    1 * 1024 * 1024, // Set BodyLimit to 1 MB
		ErrorHandler: newErrorHandler(logger),
	})

	middleware.Register(app, logger)
	route.RegisterRoutes(app, deps)

	// Penampung terakhir untuk URL yang tidak dikenal.
	app.Use(func(c *fiber.Ctx) error {
		return helper.Fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	return app
}

// newErrorHandler adalah jaring pengaman terakhir: error yang tidak
// tertangani di service berakhir di sini dengan format yang tetap konsisten.
func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		requestID := helper.RequestID(c)

		var appErr *helper.AppError

		switch {
		case errors.As(err, &appErr):

		case errors.Is(err, fiber.ErrRequestEntityTooLarge):
			appErr = &helper.AppError{
				Status: fiber.StatusRequestEntityTooLarge, Code: "PAYLOAD_TTOO_LARGE", Message: "Ukuran body melebihi batas yang diizinkan",
			}

		default:
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				appErr = &helper.AppError{
					Status: fiberErr.Code,
					Code:   "HTTP_ERROR", Message: fiberErr.Message,
				}
			} else {
				appErr = helper.Internal(err)
			}
		}

		if appErr.Status < fiber.StatusInternalServerError {
			logger.Error("request_failed",
				slog.String("request_id", requestID),
				slog.String("path", c.Path()),
				slog.String("code", appErr.Code),
				slog.Int("status", appErr.Status),
				slog.String("error", appErr.Error()))
		} else {
			logger.Warn("request_rejected",
				slog.String("request_id", requestID),
				slog.String("path", c.Path()),
				slog.String("code", appErr.Code),
				slog.Int("status", appErr.Status))
		}

		return c.Status(appErr.Status).JSON(model.ErrorResponse{
			Success:   false,
			Code:      appErr.Code,
			Message:   appErr.Message,
			Fields:    appErr.Fields,
			RequestID: requestID,
		})
		// status := fiber.StatusInternalServerError
		// message := "terjadi error pada server"

		// if e, ok := err.(*fiber.Error); ok {
		// 	status = e.Code
		// 	message = e.Message
		// }

		// logger.Error("unhandled_error",
		// 	slog.String("path", c.Path()),
		// 	slog.Int("status", status),
		// 	slog.String("error", err.Error()),
		// )

		// return helper.Fail(c, status, message)
	}
}
