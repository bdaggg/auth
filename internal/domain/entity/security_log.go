package entity

import (
	"time"

	"gorm.io/gorm"
)

type SecurityAction string

const (
	ActionBlockUser   SecurityAction = "block_user"
	ActionUnblockUser SecurityAction = "unblock_user"
	ActionFailedLogin SecurityAction = "failed_login"
	ActionSuspicious  SecurityAction = "suspicious_activity"
	ActionRoleChange  SecurityAction = "role_change"
	ActionForceLogout SecurityAction = "force_logout"
)

type SecurityLog struct {
	ID          string         `gorm:"primarykey"`
	UserID      string         `gorm:"index"`
	Action      SecurityAction `gorm:"type:varchar(50)"`
	Description string         `gorm:"type:text"`
	IP          string         `gorm:"type:varchar(45)"`
	UserAgent   string         `gorm:"type:varchar(255)"`
	Metadata    JSON           `gorm:"type:jsonb"`
	CreatedBy   string         `gorm:"type:varchar(36)"` // Eylemi gerçekleştiren kullanıcı
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
