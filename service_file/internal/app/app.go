package app

import (
	"log/slog"
	grpcapp "service-file/internal/app/grpc"
	httpapp "service-file/internal/app/http"

	http "service-file/internal/controller"

	"service-file/internal/infrastructure/postgresrepo"
	"service-file/internal/usecase"
	"service-file/pkg/config"
)

type App struct {
	HTTPSrv *httpapp.App
	GRPCSrv *grpcapp.App
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	db, err := postgresrepo.New(cfg)
	if err != nil {
		return nil, err
	}

	directoryRepo := postgresrepo.NewDirectoryRepository(db)
	fileMetadataRepo := postgresrepo.NewFileMetadataRepository(db)

	directoryUsecase := usecase.NewDirectoryUsecase(directoryRepo, logger)
	fileUsecase := usecase.NewFileUsecase(directoryRepo, fileMetadataRepo, logger)
	gRPCUsecase := usecase.NewGRPCUsecase(fileMetadataRepo, logger)

	treeHandler := http.NewTreeHandler(directoryUsecase, fileUsecase)

	httpApp := httpapp.New(logger, treeHandler, cfg)
	grpcApp := grpcapp.New(logger, gRPCUsecase, cfg.GRPCServer.Address)

	return &App{
		HTTPSrv: httpApp,
		GRPCSrv: grpcApp,
	}, nil
}
