package main

import (
	"log/slog"
	"os"
	"service-file/internal/app"
	http "service-file/internal/controller/middleware"
	"service-file/pkg/config"
	"service-file/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	log.Info("config loaded successfully",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.Address),
		slog.Any("http_server", cfg.HTTPServer),
	)

	appInstance, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := gin.Default()
	router.Use(gin.Recovery(), http.CORSMiddleware())

	setupRoutes(router, appInstance, cfg)

	port := cfg.HTTPServer.Address
	log.Info("starting server", slog.String("port", port))
	if err := router.Run(port); err != nil {
		log.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
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

func setupRoutes(router *gin.Engine, appInstance *app.App, cfg *config.Config) {
	directoriesGroup := router.Group("/directories", http.AuthMiddleware(cfg))
	{
		directoriesGroup.POST("/create", appInstance.TreeHandler.CreateDirectory)
		directoriesGroup.DELETE("", appInstance.TreeHandler.DeleteDirectory)
		directoriesGroup.POST("", appInstance.TreeHandler.GetTree)
	}

	filesGroup := router.Group("/files", http.AuthMiddleware(cfg))
	{
		filesGroup.GET("/:file_id", appInstance.TreeHandler.GetFileInfo)
		filesGroup.POST("/upload", appInstance.TreeHandler.UploadFile)
		filesGroup.DELETE("", appInstance.TreeHandler.DeleteFile)
	}
}
