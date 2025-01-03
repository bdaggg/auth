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

func (r *GormUserRepository) GetByResetToken(ctx context.Context, token string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "password_reset_token = ?", token).Error; err != nil {
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
