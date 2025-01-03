package repository

import (
	"context"
	"time"

	"auth-service/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByResetToken(ctx context.Context, token string) (*entity.User, error)
	List(ctx context.Context, offset, limit int) ([]entity.User, error)
	GetActiveCount(ctx context.Context) (int, error)
	GetBlockedCount(ctx context.Context) (int, error)
}

type SecurityRepository interface {
	CreateLog(ctx context.Context, log *entity.SecurityLog) error
	GetLogs(ctx context.Context, userID string, from, to time.Time) ([]entity.SecurityLog, error)
	GetAlerts(ctx context.Context, from, to time.Time) ([]SecurityAlert, error)
	GetSuspiciousActivities(ctx context.Context, threshold int) ([]entity.SecurityLog, error)
}

type SessionRepository interface {
	Create(ctx context.Context, sessionID string, session *entity.Session, ttl time.Duration) error
	Get(ctx context.Context, sessionID string) (*entity.Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteAllUserSessions(ctx context.Context, userID string) error
}

type AuditRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]entity.AuditLog, error)
	GetByDateRange(ctx context.Context, from, to time.Time) ([]entity.AuditLog, error)
}

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *entity.Subscription) error
	GetByUserID(ctx context.Context, userID string) ([]entity.Subscription, error)
	GetActiveSubscription(ctx context.Context, userID string, subType entity.SubscriptionType) (*entity.Subscription, error)
	Update(ctx context.Context, subscription *entity.Subscription) error
	Delete(ctx context.Context, id string) error
}

// SecurityAlert güvenlik uyarılarını temsil eder
type SecurityAlert struct {
	ID          string
	Type        string
	Severity    string
	Description string
	UserID      string
	CreatedAt   time.Time
	Metadata    map[string]interface{}
}
