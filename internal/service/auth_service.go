package service

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/oauth"
	"auth-service/internal/domain/repository"
	"auth-service/internal/domain/validator"
	"auth-service/pkg/security"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("geçersiz kimlik bilgileri")
	ErrUserExists         = errors.New("kullanıcı zaten mevcut")
)

type AuthService struct {
	userRepo        repository.UserRepository
	jwtManager      *security.JWTManager
	emailService    *EmailService
	totpService     *TOTPService
	sessionRepo     repository.SessionRepository
	auditRepo       repository.AuditRepository
	securityService *SecurityService
	oauthProvider   oauth.Provider
}

type RegisterInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

func NewAuthService(
	userRepo repository.UserRepository,
	jwtManager *security.JWTManager,
	emailService *EmailService,
	totpService *TOTPService,
	sessionRepo repository.SessionRepository,
	auditRepo repository.AuditRepository,
	securityService *SecurityService,
	oauthProvider oauth.Provider,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		jwtManager:      jwtManager,
		emailService:    emailService,
		totpService:     totpService,
		sessionRepo:     sessionRepo,
		auditRepo:       auditRepo,
		securityService: securityService,
		oauthProvider:   oauthProvider,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*entity.User, error) {
	// Email kullanımda mı kontrol et
	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Şifreyi hashle
	hashedPassword, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Yeni kullanıcı oluştur
	user := &entity.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      entity.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Kullanıcıyı kaydet
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*entity.TokenPair, error) {
	// Kullanıcıyı bul
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Şifreyi kontrol et
	if !security.CheckPassword(input.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Token pair oluştur
	return s.jwtManager.GenerateTokenPair(user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*entity.TokenPair, error) {
	// Refresh token'ı doğrula
	claims, err := s.jwtManager.ValidateToken(refreshToken, entity.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Kullanıcıyı bul
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Yeni token pair oluştur
	return s.jwtManager.GenerateTokenPair(user)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidCredentials
	}

	if !security.CheckPassword(oldPassword, user.Password) {
		return ErrInvalidCredentials
	}

	if err := validator.ValidatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) InvalidateTokens(ctx context.Context, userID string) error {
	// Redis'te kullanıcının tüm tokenlarını blacklist'e ekle
	// Bu özellik için Redis repository'si güncellenmeli
	return nil
}

func (s *AuthService) InitiatePasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return nil // Güvenlik için kullanıcı bulunamadı hatası dönmeyelim
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	user.PasswordResetToken = token
	user.TokenExpiresAt = &expiresAt

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.emailService.SendPasswordResetEmail(email, token)
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := s.userRepo.GetByResetToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil || user.TokenExpiresAt.Before(time.Now()) {
		return errors.New("geçersiz veya süresi dolmuş token")
	}

	if err := validator.ValidatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.PasswordResetToken = ""
	user.TokenExpiresAt = nil

	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) Enable2FA(ctx context.Context, userID string) (string, string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	secret, err := s.totpService.GenerateSecret()
	if err != nil {
		return "", "", err
	}

	qrCode, err := s.totpService.GenerateQRCode(user.Email, secret)
	if err != nil {
		return "", "", err
	}

	user.TOTPSecret = secret
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", "", err
	}

	return secret, qrCode, nil
}

func (s *AuthService) Verify2FA(ctx context.Context, userID, code string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !s.totpService.ValidateCode(user.TOTPSecret, code) {
		return errors.New("geçersiz 2FA kodu")
	}

	user.Is2FAEnabled = true
	return s.userRepo.Update(ctx, user)
}
