package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
	"api-students/helper"
	"api-students/route"
)

const minSecretLength = 32

func main() {
	config.LoadEnv()
	logger := config.NewLogger()

	jwtSecret := config.GetEnv("JWT_SECRET", "")
	if len(jwtSecret) < minSecretLength {
		logger.Error("JWT_SECRET tidak diisi atau terlalu pendek",
			slog.Int("minimal_karakter", minSecretLength),
		)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := database.NewPool(ctx)
	if err != nil {
		logger.Error("gagal terhubung ke database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)
	studentRepo := repository.NewStudentRepository(pool)
	roleRepo := repository.NewRoleRepository(pool)

	rawPerms, err := roleRepo.LoadPermissions(ctx)
	if err != nil {
		logger.Error("gagal memuat data permission role", slog.String("error", err.Error()))
		os.Exit(1)
	}
	perms := helper.NewPermissionSet(rawPerms)
	logger.Info("permission dimuat", slog.Any("roles", perms.KnownRoles()))

	jwtIssuer := config.GetEnv("JWT_ISSUER", "api-students")
	accessTTLMinutes := config.GetEnvInt("JWT_ACCESS_TTL_MINUTES", 15)
	refreshTTLDays := config.GetEnvInt("JWT_REFRESH_TTL_DAYS", 7)

	accessTTL := time.Duration(accessTTLMinutes) * time.Minute
	refreshTTL := time.Duration(refreshTTLDays) * 24 * time.Hour

	jwtManager := helper.NewJWTManager(jwtSecret, jwtIssuer, accessTTL)

	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager, refreshTTL, perms)
	studentService := service.NewStudentService(studentRepo, perms)

	deps := route.Dependencies{
		Pool:           pool,
		StudentService: studentService,
		AuthService:    authService,
		JWT:            jwtManager,
		Permissions:    perms,
	}

	app := config.NewApp(logger, deps)

	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("server berjalan", slog.String("port", port))

	// Graceful shutdown: tunggu Ctrl+C / SIGTERM, lalu beri waktu request yang sedang berjalan untuk selesai.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("sinyal berhenti diterima, menutup server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("gagal menutup server dengan rapi", slog.String("error", err.Error()))
	}

	logger.Info("server berhenti dengan rapi")
}
