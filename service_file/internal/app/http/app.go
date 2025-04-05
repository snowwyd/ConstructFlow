package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	controller "service-file/internal/controller"
	middleware "service-file/internal/controller/middleware"
	"service-file/pkg/config"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	log         *slog.Logger
	treeHandler *controller.TreeHandler
	cfg         *config.Config
	server      *http.Server
}

func New(log *slog.Logger, treeHandler *controller.TreeHandler, cfg *config.Config) *App {
	return &App{
		log:         log,
		treeHandler: treeHandler,
		cfg:         cfg,
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

	setupRoutes(router, a.treeHandler, a.cfg)

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

func setupRoutes(router *gin.Engine, treeHandler *controller.TreeHandler, cfg *config.Config) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	directoriesGroup := router.Group("/directories", middleware.AuthMiddleware(cfg))
	{
		directoriesGroup.POST("/create", treeHandler.CreateDirectory)
		directoriesGroup.DELETE("", treeHandler.DeleteDirectory)
		directoriesGroup.POST("", treeHandler.GetTree)
	}

	filesGroup := router.Group("/files", middleware.AuthMiddleware(cfg))
	{
		filesGroup.POST("/upload", treeHandler.UploadFile)
		filesGroup.DELETE("", treeHandler.DeleteFile)
		filesGroup.GET("/:file_id", treeHandler.GetFileInfo)
		filesGroup.PUT("/:file_id", treeHandler.UpdateFile)
		filesGroup.GET("/:file_id/download-direct", treeHandler.DownloadFileDirect)
		filesGroup.GET("/:file_id/convert/gltf", treeHandler.ConvertSTPToGLTF)

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
