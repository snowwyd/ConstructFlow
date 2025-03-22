package app

import (
	"backend/internal/controller/http"
	"backend/internal/infrastructure/postgresrepo"
	"backend/internal/usecase"
	"backend/pkg/config"
	"log/slog"
)

type App struct {
	Config          *config.Config
	Logger          *slog.Logger
	AuthHandler     *http.AuthHandler
	TreeHandler     *http.TreeHandler
	ApprovalHandler *http.ApprovalHandler
	// ... другие обработчики и use cases
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// Инициализация репозиториев
	db, err := postgresrepo.New(cfg)
	if err != nil {
		return nil, err
	}

	userRepo := postgresrepo.NewUserRepository(db)
	roleRepo := postgresrepo.NewRoleRepository(db)
	fileTreeRepo := postgresrepo.NewFileTreeRepository(db)
	approvalRepo := postgresrepo.NewApprovalRepository(db)

	// Инициализация use cases
	authUsecase := usecase.NewAuthUsecase(userRepo, roleRepo, cfg, logger)
	fileTreeUsecase := usecase.NewFileTreeUsecase(fileTreeRepo, logger)
	approvalUsecase := usecase.NewApprovalUsecase(fileTreeRepo, approvalRepo, logger)

	// Инициализация контроллеров
	authHandler := http.NewAuthHandler(authUsecase)
	treeHandler := http.NewTreeHandler(fileTreeUsecase)
	approvalHandler := http.NewApprovalHandler(approvalUsecase)

	return &App{
		Config:          cfg,
		Logger:          logger,
		AuthHandler:     authHandler,
		TreeHandler:     treeHandler,
		ApprovalHandler: approvalHandler,
	}, nil
}
