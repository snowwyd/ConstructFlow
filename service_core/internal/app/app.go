package app

import (
	"log"
	"log/slog"
	httpapp "service-core/internal/app/http"
	http "service-core/internal/controller"

	"service-core/internal/infrastructure/grpc"
	"service-core/internal/infrastructure/postgresrepo"
	"service-core/internal/usecase"
	"service-core/pkg/config"
)

type App struct {
	HTTPSrv *httpapp.App
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	db, err := postgresrepo.New(cfg)
	if err != nil {
		return nil, err
	}

	grpcClient, err := grpc.NewFileGRPCClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}

	userRepo := postgresrepo.NewUserRepository(db)
	roleRepo := postgresrepo.NewRoleRepository(db)
	approvalRepo := postgresrepo.NewApprovalRepository(db)
	workflowRepo := postgresrepo.NewWorkflowRepository(db)

	fileService := grpc.NewFileService(grpcClient)

	authUsecase := usecase.NewAuthUsecase(userRepo, roleRepo, cfg, logger)
	approvalUsecase := usecase.NewApprovalUsecase(approvalRepo, fileService, logger)
	workflowUsecase := usecase.NewWorkflowUsecase(workflowRepo, userRepo, fileService, logger)
	roleUsecase := usecase.NewRoleUsecase(roleRepo, userRepo, logger)

	authHandler := http.NewAuthHandler(authUsecase)
	fileHandler := http.NewFileHandler(approvalUsecase)
	approvalHandler := http.NewFileApprovalsHandler(approvalUsecase)
	workflowHandler := http.NewWorkflowHandler(workflowUsecase)
	roleHandler := http.NewRoleHandler(roleUsecase)

	httpApp := httpapp.New(
		logger,
		authHandler,
		fileHandler,
		approvalHandler,
		workflowHandler,
		roleHandler,
		cfg,
	)

	return &App{
		HTTPSrv: httpApp,
	}, nil
}
