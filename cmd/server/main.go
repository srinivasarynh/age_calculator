package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"github.com/srinivasarynh/age_calculator/config"
	"github.com/srinivasarynh/age_calculator/internal/handler"
	"github.com/srinivasarynh/age_calculator/internal/logger"
	"github.com/srinivasarynh/age_calculator/internal/middleware"
	"github.com/srinivasarynh/age_calculator/internal/repository"
	"github.com/srinivasarynh/age_calculator/internal/routes"
	"github.com/srinivasarynh/age_calculator/internal/service"
	"go.uber.org/zap"
)

func main() {
	zapLogger := logger.NewLogger()
	defer zapLogger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		zapLogger.Fatal("Failed to load config", zap.Error(err))
	}

	db, err := config.NewDatabase(cfg)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	zapLogger.Info("Database connection extablished")

	userRepo := repository.NewUserRepository(db, zapLogger)
	userService := service.NewUserService(userRepo, zapLogger)
	userHandler := handler.NewUserHandler(userService, zapLogger)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		AppName:      "User API v1.0",
	})

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(zapLogger))

	routes.SetupRoutes(app, userHandler)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		zapLogger.Info("Shutting down server...")
		if err := app.ShutdownWithContext(context.Background()); err != nil {
			zapLogger.Fatal("Server forced to shutdown", zap.Error(err))
		}
	}()

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	zapLogger.Info("Server starting", zap.String("address", addr))
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
