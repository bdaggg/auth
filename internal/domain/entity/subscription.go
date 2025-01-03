package entity

import (
	"time"

	"gorm.io/gorm"
)

type SubscriptionType string

const (
	BlueTick SubscriptionType = "blue_tick"
	Badge    SubscriptionType = "badge"
)

type Subscription struct {
	ID        string           `gorm:"primarykey"`
	UserID    string           `gorm:"not null"`
	Type      SubscriptionType `gorm:"type:varchar(20);not null"`
	StartDate time.Time        `gorm:"not null"`
	EndDate   time.Time        `gorm:"not null"`
	IsActive  bool             `gorm:"default:true"`
	PaymentID string           `gorm:"uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// İlişkiler
	User User `gorm:"foreignKey:UserID"`
}
