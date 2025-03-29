package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	controller "service-core/internal/controller"
	middleware "service-core/internal/controller/middleware"
	"service-core/pkg/config"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	log             *slog.Logger
	authHandler     *controller.AuthHandler
	approvalHandler *controller.ApprovalHandler
	cfg             *config.Config
	server          *http.Server
}

func New(log *slog.Logger, authHandler *controller.AuthHandler, approvalHandler *controller.ApprovalHandler, cfg *config.Config) *App {
	return &App{
		log:             log,
		authHandler:     authHandler,
		approvalHandler: approvalHandler,
		cfg:             cfg,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("operation", op),
		slog.String("port", a.cfg.HTTPServer.Address),
	)

	router := gin.Default()
	router.Use(gin.Recovery(), middleware.CORSMiddleware())

	setupRoutes(router, a.authHandler, a.approvalHandler, a.cfg)

	port := a.cfg.HTTPServer.Address
	log.Info("starting HTTP server", slog.String("port", port))

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	a.server = srv

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func setupRoutes(router *gin.Engine, authHandler *controller.AuthHandler, approvalHandler *controller.ApprovalHandler, cfg *config.Config) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.RegisterUser)
		authGroup.POST("/role", authHandler.RegisterRole)

		authGroup.GET("/me", middleware.AuthMiddleware(cfg), authHandler.GetCurrentUser)
	}

	filesGroup := router.Group("/files", middleware.AuthMiddleware(cfg))
	{
		filesGroup.PUT("/:file_id/approve", approvalHandler.ApproveFile)
	}

	approvalsGroup := router.Group("/file-approvals", middleware.AuthMiddleware(cfg))
	{
		approvalsGroup.GET("", approvalHandler.GetApprovalsByUser)
		approvalsGroup.PUT("/:approval_id/sign", approvalHandler.SignApproval)
		approvalsGroup.PUT("/:approval_id/annotate", approvalHandler.AnnotateApproval)
		approvalsGroup.PUT("/:approval_id/finalize", approvalHandler.FinalizeApproval)
	}
}

func (a *App) Stop() {
	a.log.Info("http server is stopping")
	if a.server == nil {
		a.log.Info("http server is not running")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("failed to gracefully shutdown http server", slog.String("error", err.Error()))
	} else {
		a.log.Info("http server stopped gracefully")
	}
}
