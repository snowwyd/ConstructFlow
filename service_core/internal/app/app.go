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

	grpcClient, err := grpc.NewFileGRPCClient(cfg.GRPCAddress)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}

	userRepo := postgresrepo.NewUserRepository(db)
	roleRepo := postgresrepo.NewRoleRepository(db)
	approvalRepo := postgresrepo.NewApprovalRepository(db)
	fileService := grpc.NewFileService(grpcClient)

	authUsecase := usecase.NewAuthUsecase(userRepo, roleRepo, cfg, logger)
	approvalUsecase := usecase.NewApprovalUsecase(approvalRepo, fileService, logger)

	authHandler := http.NewAuthHandler(authUsecase)
	approvalHandler := http.NewApprovalHandler(approvalUsecase)

	httpApp := httpapp.New(logger, authHandler, approvalHandler, cfg)

	return &App{
		HTTPSrv: httpApp,
	}, nil
}
