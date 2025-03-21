package postgresrepo

import (
	"backend/internal/domain"
	"context"

	"gorm.io/gorm"
)

type ApprovalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *Database) *ApprovalRepository {
	return &ApprovalRepository{db: db.db}
}

func (r *ApprovalRepository) CreateApproval(ctx context.Context, approval *domain.Approval, tx *gorm.DB) error {
	return tx.WithContext(ctx).Create(approval).Error
}
