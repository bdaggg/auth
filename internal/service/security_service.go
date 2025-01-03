package service

import (
	"context"
	"time"

	"auth-service/internal/domain/entity"
)

type SecurityService struct {
	userRepo     repository.UserRepository
	securityRepo repository.SecurityRepository
	sessionRepo  repository.SessionRepository
	monitoring   *MonitoringService
}

func NewSecurityService(
	userRepo repository.UserRepository,
	securityRepo repository.SecurityRepository,
	sessionRepo repository.SessionRepository,
	monitoring *MonitoringService,
) *SecurityService {
	return &SecurityService{
		userRepo:     userRepo,
		securityRepo: securityRepo,
		sessionRepo:  sessionRepo,
		monitoring:   monitoring,
	}
}

func (s *SecurityService) BlockUser(ctx context.Context, userID, blockedBy, reason string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.IsActive = false
	user.BlockedAt = &now
	user.BlockedBy = blockedBy
	user.BlockReason = reason

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Kullanıcının tüm oturumlarını sonlandır
	if err := s.sessionRepo.DeleteAllUserSessions(ctx, userID); err != nil {
		return err
	}

	// Güvenlik logu oluştur
	secLog := &entity.SecurityLog{
		UserID:      userID,
		Action:      entity.ActionBlockUser,
		Description: reason,
		CreatedBy:   blockedBy,
		CreatedAt:   now,
	}

	return s.securityRepo.CreateLog(ctx, secLog)
}

func (s *SecurityService) UnblockUser(ctx context.Context, userID, unblockBy string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.IsActive = true
	user.BlockedAt = nil
	user.BlockedBy = ""
	user.BlockReason = ""

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	secLog := &entity.SecurityLog{
		UserID:    userID,
		Action:    entity.ActionUnblockUser,
		CreatedBy: unblockBy,
		CreatedAt: time.Now(),
	}

	return s.securityRepo.CreateLog(ctx, secLog)
}

func (s *SecurityService) CheckSuspiciousActivity(ctx context.Context, userID string) error {
	// Son başarısız giriş denemelerini kontrol et
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.FailedLoginAttempts >= 5 {
		// Hesabı otomatik olarak bloke et
		return s.BlockUser(ctx, userID, "system", "Too many failed login attempts")
	}

	return nil
}
