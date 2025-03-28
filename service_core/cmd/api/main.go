package main

import (
	"log/slog"
	"os"
	"service-core/internal/app"
	http "service-core/internal/controller/middleware"
	"service-core/pkg/config"
	"service-core/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "service-core/docs"
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

func setupRoutes(router *gin.Engine, appInstance *app.App, cfg *config.Config) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Маршруты аутентификации
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", appInstance.AuthHandler.Login)
		authGroup.POST("/register", appInstance.AuthHandler.RegisterUser)
		authGroup.POST("/role", appInstance.AuthHandler.RegisterRole)

		// Защищен middleware
		authGroup.GET("/me", http.AuthMiddleware(cfg), appInstance.AuthHandler.GetCurrentUser)
	}

	filesGroup := router.Group("/files", http.AuthMiddleware(cfg))
	{
		filesGroup.PUT("/:file_id/approve", appInstance.ApprovalHandler.ApproveFile)
	}

	approvalsGroup := router.Group("/file-approvals", http.AuthMiddleware(cfg))
	{
		approvalsGroup.GET("", appInstance.ApprovalHandler.GetApprovalsByUser)
		approvalsGroup.PUT("/:approval_id/sign", appInstance.ApprovalHandler.SignApproval)
		approvalsGroup.PUT("/:approval_id/annotate", appInstance.ApprovalHandler.AnnotateApproval)
		approvalsGroup.PUT("/:approval_id/finalize", appInstance.ApprovalHandler.FinalizeApproval)
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
