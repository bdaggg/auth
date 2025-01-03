package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID string    `json:"user_id"`
	Role   Role      `json:"role"`
	Type   TokenType `json:"token_type"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
