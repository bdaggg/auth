package entity

import "time"

type Session struct {
	UserID    string
	ExpiresAt time.Time
	Data      map[string]interface{}
}
