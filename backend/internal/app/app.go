package app

import (
	"backend/internal/controller/http"
	"backend/internal/infrastructure/postgresrepo"
	"backend/internal/usecase"
	"backend/pkg/config"
	"log/slog"
)

type App struct {
	Config      *config.Config
	Logger      *slog.Logger
	AuthHandler *http.AuthHandler
	// ... другие обработчики и use cases
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// Инициализация репозиториев
	db, err := postgresrepo.New(cfg)
	if err != nil {
		return nil, err
	}

	userRepo := postgresrepo.NewUserRepository(db)

	// TODO: Инициализация use cases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg, logger)

	// TODO: Инициализация контроллеров
	authHandler := http.NewAuthHandler(authUsecase, cfg)

	return &App{
		Config:      cfg,
		Logger:      logger,
		AuthHandler: authHandler,
	}, nil
}
