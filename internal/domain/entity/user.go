package entity

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser    Role = "user"
	RoleAdmin   Role = "admin"
	RolePremium Role = "premium"
)

type User struct {
	ID                     string `gorm:"primarykey"`
	Email                  string `gorm:"uniqueIndex;not null"`
	Password               string `gorm:"not null"`
	Role                   Role   `gorm:"type:varchar(20);default:'user'"`
	IsVerified             bool   `gorm:"default:false"`
	HasBlueTick            bool   `gorm:"default:false"`
	Is2FAEnabled           bool   `gorm:"default:false"`
	TOTPSecret             string `gorm:"type:varchar(32)"`
	EmailVerificationToken string `gorm:"type:varchar(100)"`
	PasswordResetToken     string `gorm:"type:varchar(100)"`
	TokenExpiresAt         *time.Time
	IsActive               bool `gorm:"default:true"`
	BlockedAt              *time.Time
	BlockedBy              string `gorm:"type:varchar(36)"`
	BlockReason            string `gorm:"type:text"`
	LastLoginAt            *time.Time
	LastLoginIP            string `gorm:"type:varchar(45)"`
	FailedLoginAttempts    int    `gorm:"default:0"`
	LastFailedLoginAt      *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`

	// İlişkiler
	Subscriptions []Subscription `gorm:"foreignKey:UserID"`
	AuditLogs     []AuditLog     `gorm:"foreignKey:UserID"`
	SecurityLogs  []SecurityLog  `gorm:"foreignKey:UserID"`
}
