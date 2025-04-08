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
	log                  *slog.Logger
	authHandler          *controller.AuthHandler
	fileHandler          *controller.FileHandler
	fileApprovalsHandler *controller.FileApprovalsHandler
	workflowHandler      *controller.WorkflowlHandler
	roleHandler          *controller.RoleHandler
	userHandler          *controller.UserHandler
	cfg                  *config.Config
	server               *http.Server
}

func New(
	log *slog.Logger,
	authHandler *controller.AuthHandler,
	fileHandler *controller.FileHandler,
	fileApprovalsHandler *controller.FileApprovalsHandler,
	workflowHandler *controller.WorkflowlHandler,
	roleHandler *controller.RoleHandler,
	userHandler *controller.UserHandler,
	cfg *config.Config,
) *App {
	return &App{
		log:                  log,
		authHandler:          authHandler,
		fileHandler:          fileHandler,
		fileApprovalsHandler: fileApprovalsHandler,
		workflowHandler:      workflowHandler,
		roleHandler:          roleHandler,
		userHandler:          userHandler,
		cfg:                  cfg,
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

	setupRoutes(
		router,
		a.authHandler,
		a.fileHandler,
		a.fileApprovalsHandler,
		a.workflowHandler,
		a.roleHandler,
		a.userHandler,
		a.cfg,
	)

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

func setupRoutes(
	router *gin.Engine,
	authHandler *controller.AuthHandler,
	fileHandler *controller.FileHandler,
	fileApprovalsHandler *controller.FileApprovalsHandler,
	workflowHandler *controller.WorkflowlHandler,
	roleHandler *controller.RoleHandler,
	userHandler *controller.UserHandler,
	cfg *config.Config,
) {
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/me", middleware.AuthMiddleware(cfg), authHandler.GetCurrentUser)
	}

	filesGroup := router.Group("/files", middleware.AuthMiddleware(cfg))
	{
		filesGroup.PUT("/:file_id/approve", fileHandler.ApproveFile)
	}

	approvalsGroup := router.Group("/file-approvals", middleware.AuthMiddleware(cfg))
	{
		approvalsGroup.GET("", fileApprovalsHandler.GetApprovalsByUser)
		approvalsGroup.PUT("/:approval_id/sign", fileApprovalsHandler.SignApproval)
		approvalsGroup.PUT("/:approval_id/annotate", fileApprovalsHandler.AnnotateApproval)
		approvalsGroup.PUT("/:approval_id/finalize", fileApprovalsHandler.FinalizeApproval)
	}

	adminGroup := router.Group("/admin", middleware.AuthMiddleware(cfg))
	{
		workflowsGroup := adminGroup.Group("/workflows")
		{
			workflowsGroup.GET("", workflowHandler.GetWorkflows)
			workflowsGroup.POST("", workflowHandler.CreateWorkflow)
			workflowsGroup.DELETE("", workflowHandler.DeleteWorkflow)
			// TODO: get workflow by id
			workflowsGroup.PUT("/:workflow_id", workflowHandler.UpdateWorkflow)
		}

		rolesGroup := adminGroup.Group("/roles")
		{
			rolesGroup.GET("", roleHandler.GetRoles)
			rolesGroup.GET("/:role_id", roleHandler.GetRole)
			rolesGroup.POST("", roleHandler.RegisterRole)
			rolesGroup.PUT("/:role_id", roleHandler.UpdateRole)
			rolesGroup.DELETE("", roleHandler.DeleteRole)
		}

		usersGroup := adminGroup.Group("/users")
		{
			usersGroup.GET("", userHandler.GetUsers)
			usersGroup.POST("/register", userHandler.RegisterUser)
			usersGroup.PUT("/:user_id", userHandler.UpdateUser)
			usersGroup.DELETE("", userHandler.DeleteUser)
		}
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
