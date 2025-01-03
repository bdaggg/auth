package repository

import (
	"context"
	"time"

	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/repository"

	"gorm.io/gorm"
)

type GormAuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) repository.AuditRepository {
	return &GormAuditRepository{db: db}
}

func (r *GormAuditRepository) Create(ctx context.Context, log *entity.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *GormAuditRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&logs).Error
	return logs, err
}

func (r *GormAuditRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.WithContext(ctx).Where("created_at BETWEEN ? AND ?", from, to).Find(&logs).Error
	return logs, err
}
