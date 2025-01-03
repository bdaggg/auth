package repository

import (
	"context"
	"errors"

	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/repository"

	"gorm.io/gorm"
)

type GormSubscriptionRepository struct {
	db *gorm.DB
}

func NewGormSubscriptionRepository(db *gorm.DB) repository.SubscriptionRepository {
	return &GormSubscriptionRepository{db: db}
}

func (r *GormSubscriptionRepository) Create(ctx context.Context, subscription *entity.Subscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

func (r *GormSubscriptionRepository) GetByUserID(ctx context.Context, userID string) ([]entity.Subscription, error) {
	var subscriptions []entity.Subscription
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *GormSubscriptionRepository) GetActiveSubscription(ctx context.Context, userID string, subType entity.SubscriptionType) (*entity.Subscription, error) {
	var subscription entity.Subscription
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND is_active = ? AND end_date > NOW()", userID, subType, true).
		First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &subscription, nil
}

func (r *GormSubscriptionRepository) Update(ctx context.Context, subscription *entity.Subscription) error {
	return r.db.WithContext(ctx).Save(subscription).Error
}

func (r *GormSubscriptionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Subscription{}, "id = ?", id).Error
}
