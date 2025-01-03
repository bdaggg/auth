package repository

import (
	"context"
	"errors"

	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/repository"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repository.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *GormUserRepository) List(ctx context.Context, offset, limit int) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

func (r *GormUserRepository) GetActiveCount(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("is_active = ?", true).Count(&count).Error
	return int(count), err
}

func (r *GormUserRepository) GetBlockedCount(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("is_active = ?", false).Count(&count).Error
	return int(count), err
}

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
