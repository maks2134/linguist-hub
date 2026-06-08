package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	swagger "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	_ "github.com/makskozlov/linguist-hub/docs"
	"github.com/makskozlov/linguist-hub/internal/config"
	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/internal/containers"
	"github.com/makskozlov/linguist-hub/internal/models"
	"github.com/makskozlov/linguist-hub/internal/services"
	"github.com/makskozlov/linguist-hub/pkg/middleware"
	"github.com/makskozlov/linguist-hub/pkg/postgres"
	"gorm.io/gorm"
)

type App struct {
	fiber         *fiber.App
	cfg           *config.Config
	db            *gorm.DB
	geminiService *services.GeminiService
	logger        *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	db, err := postgres.NewDB(cfg.DBDsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.LanguageProfile{},
		&models.StudyPlan{},
		&models.ChatMessage{},
		&models.VoiceLesson{},
	); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed")

	geminiService, err := services.NewGeminiService(cfg.GeminiAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini service: %w", err)
	}

	logger.Info("Gemini service initialized")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			logger.Error("fiber error", slog.Any("error", err))
			return c.Status(constants.StatusInternalError).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "Internal server error",
				},
			})
		},
	})

	app.Use(middleware.RequestID())
	app.Use(middleware.Recover(logger))
	app.Use(middleware.Logger(logger))
	app.Use(middleware.CORS([]string{cfg.AllowedOrigins}))

	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	container := containers.NewContainer(db, geminiService, cfg.JWTSecret, cfg.FeatureServerSTT)

	api := app.Group("/api")
	authMiddleware := middleware.JWTAuth(cfg.JWTSecret)

	container.RegisterRoutes(api, authMiddleware)

	logger.Info("Routes registered")

	return &App{
		fiber:         app,
		cfg:           cfg,
		db:            db,
		geminiService: geminiService,
		logger:        logger,
	}, nil
}

func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		a.logger.Info("Shutting down server...")
		if err := a.Shutdown(); err != nil {
			a.logger.Error("Error during shutdown", slog.Any("error", err))
		}
	}()

	addr := fmt.Sprintf(":%s", a.cfg.Port)
	a.logger.Info("Starting server", slog.String("addr", addr))

	if err := a.fiber.Listen(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (a *App) Shutdown() error {
	if err := a.fiber.Shutdown(); err != nil {
		a.logger.Error("error shutting down fiber", slog.Any("error", err))
	}

	if a.geminiService != nil {
		if err := a.geminiService.Close(); err != nil {
			a.logger.Error("error closing gemini service", slog.Any("error", err))
		}
	}

	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				a.logger.Error("error closing database", slog.Any("error", err))
			}
		}
	}

	a.logger.Info("Server shut down gracefully")
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.logger.Error("failed to run application", slog.Any("error", err))
		os.Exit(1)
	}
}

func Run() error {
	cfg := config.MustLoad()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	app, err := New(cfg, logger)
	if err != nil {
		return err
	}

	return app.Run()
}
