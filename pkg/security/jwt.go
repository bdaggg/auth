package security

import (
	"fmt"
	"time"

	"auth-service/internal/domain/entity"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	Issuer             string
}

type JWTManager struct {
	config JWTConfig
}

func NewJWTManager(config JWTConfig) *JWTManager {
	return &JWTManager{config: config}
}

func (m *JWTManager) GenerateTokenPair(user *entity.User) (*entity.TokenPair, error) {
	// Access Token oluştur
	accessToken, err := m.generateToken(user, entity.AccessToken, m.config.AccessTokenSecret, m.config.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("access token oluşturulamadı: %w", err)
	}

	// Refresh Token oluştur
	refreshToken, err := m.generateToken(user, entity.RefreshToken, m.config.RefreshTokenSecret, m.config.RefreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("refresh token oluşturulamadı: %w", err)
	}

	return &entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *JWTManager) generateToken(user *entity.User, tokenType entity.TokenType, secret string, ttl time.Duration) (string, error) {
	claims := &entity.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    m.config.Issuer,
			Subject:   user.ID,
		},
		UserID: user.ID,
		Role:   user.Role,
		Type:   tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (m *JWTManager) ValidateToken(tokenString string, tokenType entity.TokenType) (*entity.TokenClaims, error) {
	var secret string
	switch tokenType {
	case entity.AccessToken:
		secret = m.config.AccessTokenSecret
	case entity.RefreshToken:
		secret = m.config.RefreshTokenSecret
	default:
		return nil, fmt.Errorf("geçersiz token tipi")
	}

	token, err := jwt.ParseWithClaims(tokenString, &entity.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("beklenmeyen imza metodu: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token ayrıştırılamadı: %w", err)
	}

	if claims, ok := token.Claims.(*entity.TokenClaims); ok && token.Valid {
		if claims.Type != tokenType {
			return nil, fmt.Errorf("token tipi uyuşmuyor")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("geçersiz token")
}
