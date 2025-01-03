package repository

import (
	"context"
	"time"

	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/repository"

	"gorm.io/gorm"
)

type GormSecurityRepository struct {
	db *gorm.DB
}

func NewSecurityRepository(db *gorm.DB) repository.SecurityRepository {
	return &GormSecurityRepository{db: db}
}

func (r *GormSecurityRepository) CreateLog(ctx context.Context, log *entity.SecurityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *GormSecurityRepository) GetLogs(ctx context.Context, userID string, from, to time.Time) ([]entity.SecurityLog, error) {
	var logs []entity.SecurityLog
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, from, to).
		Find(&logs).Error
	return logs, err
}

func (r *GormSecurityRepository) GetAlerts(ctx context.Context, from, to time.Time) ([]repository.SecurityAlert, error) {
	var alerts []repository.SecurityAlert
	err := r.db.WithContext(ctx).
		Model(&entity.SecurityLog{}).
		Where("created_at BETWEEN ? AND ? AND severity = ?", from, to, "high").
		Find(&alerts).Error
	return alerts, err
}

func (r *GormSecurityRepository) GetSuspiciousActivities(ctx context.Context, threshold int) ([]entity.SecurityLog, error) {
	var logs []entity.SecurityLog
	err := r.db.WithContext(ctx).
		Where("action = ? AND created_at > ?", entity.ActionSuspicious, time.Now().Add(-24*time.Hour)).
		Having("COUNT(*) >= ?", threshold).
		Group("user_id").
		Find(&logs).Error
	return logs, err
}
