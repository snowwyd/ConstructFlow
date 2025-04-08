package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/logger/slogger"
)

type WorkflowUsecase struct {
	workflowRepo interfaces.WorkflowRepository
	userRepo     interfaces.UserRepository
	fileService  interfaces.FileService
	log          *slog.Logger
}

func NewWorkflowUsecase(
	workflowRepo interfaces.WorkflowRepository,
	userRepo interfaces.UserRepository,
	fileService interfaces.FileService,
	log *slog.Logger,
) *WorkflowUsecase {
	return &WorkflowUsecase{
		workflowRepo: workflowRepo,
		userRepo:     userRepo,
		fileService:  fileService,
		log:          log,
	}
}

func (workflowUsecase *WorkflowUsecase) GetWorkflows(ctx context.Context, userID uint) ([]domain.WorkflowResponse, error) {
	const op = "usecase.workflow.GetWorkflows"

	log := workflowUsecase.log.With(slog.String("op", op), slog.Any("user_id", userID))
	log.Info("getting workflows")

	log.Debug("checking if user is admin")
	if err := workflowUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("getting worklows from database")
	workflows, err := workflowUsecase.workflowRepo.GetWorkflows(ctx)
	if err != nil {
		// TODO: custom errors?
		log.Error("failed to get workflows", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("workflows got succesfully")
	return workflows, nil
}

func (workflowUsecase *WorkflowUsecase) CreateWorkflow(ctx context.Context, name string, stages []domain.WorkflowStage, userID uint) error {
	const op = "usecase.workflow.CreateWorkflow"

	log := workflowUsecase.log.With(slog.String("op", op))
	log.Info("creating workfow")

	log.Debug("checking if user is admin")
	if err := workflowUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking users from request")
	if err := workflowUsecase.checkUsers(ctx, stages); err != nil {
		log.Error("failed to validate users", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("putting workflow into db")
	if err := workflowUsecase.workflowRepo.CreateWorkflow(ctx, name, stages); err != nil {
		// TODO: custom errors?
		log.Error("failed to create workflow", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("workflow created successfully")
	return nil
}

func (workflowUsecase *WorkflowUsecase) UpdateWorkflow(ctx context.Context, workflowID uint, name string, stages []domain.WorkflowStage, userID uint) error {
	const op = "usecase.workflow.UpdateWorkflow"

	log := workflowUsecase.log.With(slog.String("op", op))
	log.Info("updating workfow")

	log.Debug("checking if user is admin")
	if err := workflowUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking workflow existence")
	if err := workflowUsecase.checkWorkflow(ctx, workflowID); err != nil {
		log.Error("failed workflow existence check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking users from request")
	if err := workflowUsecase.checkUsers(ctx, stages); err != nil {
		log.Error("failed to validate users", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("updating workflow")
	if err := workflowUsecase.workflowRepo.UpdateWorkflow(ctx, workflowID, name, stages); err != nil {
		// TODO: custom errors?
		log.Error("failed to update workflow", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("workflow updated successfully")
	return nil
}

func (workflowUsecase *WorkflowUsecase) DeleteWorkflow(ctx context.Context, workflowID uint, userID uint) error {
	const op = "usecase.workflow.DeleteWorkflow"

	log := workflowUsecase.log.With(slog.String("op", op))
	log.Info("deleting workfow")

	log.Debug("checking if user is admin")
	if err := workflowUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking workflow existence")
	if err := workflowUsecase.checkWorkflow(ctx, workflowID); err != nil {
		log.Error("failed workflow existence check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking if this workflow in use")
	exists, err := workflowUsecase.fileService.CheckWorkflow(ctx, workflowID)
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
	if err := workflowUsecase.workflowRepo.DeleteWorkflow(ctx, workflowID); err != nil {
		// TODO: custom errors
		log.Error("failed to delete workflow", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (workflowUsecase *WorkflowUsecase) AssignWorkflow(ctx context.Context, workflowID uint, directoryIDs []uint, userID uint) error {
	const op = "usecase.workflow.AssignWorkflow"

	log := workflowUsecase.log.With(slog.String("op", op))
	log.Info("assigning workfow")

	log.Debug("checking if user is admin")
	if err := workflowUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking workflow existence")
	if err := workflowUsecase.checkWorkflow(ctx, workflowID); err != nil {
		log.Error("failed workflow existence check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	directories := make([]uint32, len(directoryIDs))
	for i, value := range directoryIDs {
		directories[i] = uint32(value)
	}

	log.Debug("deleting workflow")
	if err := workflowUsecase.fileService.AssignWorkflow(ctx, workflowID, directories); err != nil {
		// TODO: custom errors
		log.Error("failed to delete workflow", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (workflowUsecase *WorkflowUsecase) checkAdmin(ctx context.Context, userID uint) error {
	role, err := workflowUsecase.userRepo.GetUserRole(ctx, userID)
	if err != nil {
		return err
	}
	if role != "admin" {
		return domain.ErrAccessDenied
	}
	return nil
}

func (workflowUsecase *WorkflowUsecase) checkUsers(ctx context.Context, stages []domain.WorkflowStage) error {
	userMap := make(map[uint]struct{})
	for _, stage := range stages {
		userMap[stage.UserID] = struct{}{}
	}

	userIDs := make([]uint, 0, len(userMap))
	for userID := range userMap {
		userIDs = append(userIDs, userID)
	}

	allExists, err := workflowUsecase.userRepo.CheckUsersExist(ctx, userIDs)
	if err != nil {
		return err
	}
	if !allExists {
		return domain.ErrUserNotFound
	}

	return nil
}

func (workflowUsecase *WorkflowUsecase) checkWorkflow(ctx context.Context, workflowID uint) error {
	exists, err := workflowUsecase.workflowRepo.CheckWorkflow(ctx, workflowID)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrWorkflowNotFound
	}
	return nil
}
