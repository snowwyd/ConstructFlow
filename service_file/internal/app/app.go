package app

import (
	"log/slog"
	http "service-file/internal/controller"
	grpcHandler "service-file/internal/infrastructure/grpc"

	"service-file/internal/infrastructure/postgresrepo"
	"service-file/internal/usecase"
	"service-file/pkg/config"
)

type App struct {
	Config      *config.Config
	Logger      *slog.Logger
	TreeHandler *http.TreeHandler
	GRPCServer  *grpcHandler.GRPCServer // Добавлено поле для gRPC-сервера

	// ... другие обработчики и use cases
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// Инициализация репозиториев
	db, err := postgresrepo.New(cfg)
	if err != nil {
		return nil, err
	}

	fileTreeRepo := postgresrepo.NewFileTreeRepository(db)

	// Инициализация use cases
	fileTreeUsecase := usecase.NewFileTreeUsecase(fileTreeRepo, logger)

	// Инициализация контроллеров
	treeHandler := http.NewTreeHandler(fileTreeUsecase)
	grpcServer := grpcHandler.NewGRPCServer(fileTreeUsecase)

	// Инициализация gRPC-сервера
	return &App{
		Config:      cfg,
		Logger:      logger,
		TreeHandler: treeHandler,
		GRPCServer:  grpcServer,
	}, nil
}
