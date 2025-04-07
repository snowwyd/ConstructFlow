package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/logger/slogger"
)

type WorkflowlUsecase struct {
	workflowRepo interfaces.WorkflowRepository
	userRepo     interfaces.UserRepository
	fileService  interfaces.FileService
	log          *slog.Logger
}

func NewWorkflowUsecase(workflowRepo interfaces.WorkflowRepository, userRepo interfaces.UserRepository, fileService interfaces.FileService, log *slog.Logger) *WorkflowlUsecase {
	return &WorkflowlUsecase{
		workflowRepo: workflowRepo,
		userRepo:     userRepo,
		fileService:  fileService,
		log:          log,
	}
}

func (workflowlUsecase *WorkflowlUsecase) GetWorkflows(ctx context.Context, userID uint) ([]domain.WorkflowResponse, error) {
	const op = "usecase.workflow.GetWorkflows"

	log := workflowlUsecase.log.With(slog.String("op", op), slog.Any("user_id", userID))
	log.Info("getting workflows")

	log.Debug("checking if user is admin")
	if err := workflowlUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("getting worklows from database")
	workflows, err := workflowlUsecase.workflowRepo.GetWorkflows(ctx)
	if err != nil {
		// TODO: custom errors?
		log.Error("failed to get workflows", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("workflows got succesfully")
	return workflows, nil
}

func (workflowlUsecase *WorkflowlUsecase) CreateWorkflow(ctx context.Context, name string, stages []domain.WorkflowStage, userID uint) error {
	const op = "usecase.workflow.CreateWorkflow"

	log := workflowlUsecase.log.With(slog.String("op", op))
	log.Info("creating workfow")

	log.Debug("checking if user is admin")
	if err := workflowlUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking users from request")
	if err := workflowlUsecase.checkUsers(ctx, stages); err != nil {
		log.Error("failed to validate users", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("putting workflow into db")
	if err := workflowlUsecase.workflowRepo.CreateWorkflow(ctx, name, stages); err != nil {
		// TODO: custom errors?
		log.Error("failed to create workflow", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("workflow created successfully")
	return nil
}

func (workflowlUsecase *WorkflowlUsecase) UpdateWorkflow(ctx context.Context, workflowID uint, name string, stages []domain.WorkflowStage, userID uint) error {
	const op = "usecase.workflow.UpdateWorkflow"

	log := workflowlUsecase.log.With(slog.String("op", op))
	log.Info("updating workfow")

	log.Debug("checking if user is admin")
	if err := workflowlUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking workflow existence")
	if err := workflowlUsecase.checkWorkflow(ctx, workflowID); err != nil {
		log.Error("failed workflow existence check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking users from request")
	if err := workflowlUsecase.checkUsers(ctx, stages); err != nil {
		log.Error("failed to validate users", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("updating workflow")
	if err := workflowlUsecase.workflowRepo.UpdateWorkflow(ctx, workflowID, name, stages); err != nil {
		// TODO: custom errors?
		log.Error("failed to update workflow", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("workflow updated successfully")
	return nil
}

func (workflowlUsecase *WorkflowlUsecase) DeleteWorkflow(ctx context.Context, workflowID uint, userID uint) error {
	const op = "usecase.workflow.DeleteWorkflow"

	log := workflowlUsecase.log.With(slog.String("op", op))
	log.Info("deleting workfow")

	log.Debug("checking if user is admin")
	if err := workflowlUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking workflow existence")
	if err := workflowlUsecase.checkWorkflow(ctx, workflowID); err != nil {
		log.Error("failed workflow existence check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking if this workflow in use")
	exists, err := workflowlUsecase.fileService.CheckWorkflow(ctx, workflowID)
	if err != nil {
		// TODO: custom errors
		log.Error("failed workflow check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if exists {
		log.Error("workflow is in use", slogger.Err(domain.ErrWorkflowInUse))
		return fmt.Errorf("%s: %w", op, domain.ErrWorkflowInUse)
	}

	log.Debug("deleting workflow")
	if err := workflowlUsecase.workflowRepo.DeleteWorkflow(ctx, workflowID); err != nil {
		// TODO: custom errors
		log.Error("failed to delete workflow", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (workflowlUsecase *WorkflowlUsecase) checkAdmin(ctx context.Context, userID uint) error {
	role, err := workflowlUsecase.userRepo.GetUserRole(ctx, userID)
	if err != nil {
		return err
	}
	if role != "admin" {
		return domain.ErrAccessDenied
	}
	return nil
}

func (workflowlUsecase *WorkflowlUsecase) checkUsers(ctx context.Context, stages []domain.WorkflowStage) error {
	userMap := make(map[uint]struct{})
	for _, stage := range stages {
		userMap[stage.UserID] = struct{}{}
	}

	userIDs := make([]uint, 0, len(userMap))
	for userID := range userMap {
		userIDs = append(userIDs, userID)
	}

	allExists, err := workflowlUsecase.userRepo.CheckUsersExist(ctx, userIDs)
	if err != nil {
		return err
	}
	if !allExists {
		return domain.ErrUserNotFound
	}

	return nil
}

func (workflowlUsecase *WorkflowlUsecase) checkWorkflow(ctx context.Context, workflowID uint) error {
	exists, err := workflowlUsecase.workflowRepo.CheckWorkflow(ctx, workflowID)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrWorkflowNotFound
	}
	return nil
}
