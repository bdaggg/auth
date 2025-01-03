package repository

import (
	"context"

	"auth-service/internal/domain/entity"

	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, log *entity.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AuditRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}
