package main

import (
	"backend/internal/app"
	http "backend/internal/controller/http/middleware"
	"backend/pkg/config"
	"backend/pkg/logger"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	// log.Info("config loaded successfully",
	// 	slog.String("env", cfg.Env),
	// 	slog.String("address", cfg.Address),
	// 	slog.Any("http_server", cfg.HTTPServer),
	// )
	log.Info("config loaded successfully", slog.Any("config", cfg))

	appInstance, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery())

	setupRoutes(router, appInstance, cfg)

	port := cfg.HTTPServer.Address
	log.Info("starting server", slog.String("port", port))
	if err := router.Run(port); err != nil {
		log.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

}

func setupRoutes(router *gin.Engine, appInstance *app.App, cfg *config.Config) {
	// Группа API v1
	api := router.Group("/api/v1")

	// Маршруты аутентификации
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", appInstance.AuthHandler.Login)
		authGroup.POST("/register", appInstance.AuthHandler.RegisterUser)
		authGroup.POST("/role", appInstance.AuthHandler.RegisterRole)

		// Защищен middleware
		authGroup.GET("/me", http.AuthMiddleware(cfg), appInstance.AuthHandler.GetCurrentUser)
	}
}

func setupLogger() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
