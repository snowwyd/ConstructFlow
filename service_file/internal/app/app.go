package app

import (
	"log/slog"
	http "service-file/internal/controller"
	"service-file/internal/infrastructure/postgresrepo"
	"service-file/internal/usecase"
	"service-file/pkg/config"
)

type App struct {
	Config      *config.Config
	Logger      *slog.Logger
	TreeHandler *http.TreeHandler
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

	return &App{
		Config:      cfg,
		Logger:      logger,
		TreeHandler: treeHandler,
	}, nil
}
