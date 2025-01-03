package entity

import (
	"time"

	"gorm.io/gorm"
)

type AuditAction string

const (
	ActionLogin          AuditAction = "login"
	ActionLogout         AuditAction = "logout"
	ActionPasswordChange AuditAction = "password_change"
	ActionEmailVerify    AuditAction = "email_verify"
	Action2FAEnable      AuditAction = "2fa_enable"
	Action2FADisable     AuditAction = "2fa_disable"
)

type AuditLog struct {
	ID        string      `gorm:"primarykey"`
	UserID    string      `gorm:"index"`
	Action    AuditAction `gorm:"type:varchar(50)"`
	IP        string      `gorm:"type:varchar(45)"`
	UserAgent string      `gorm:"type:varchar(255)"`
	Status    bool
	Details   string `gorm:"type:text"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
