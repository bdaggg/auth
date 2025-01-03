package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
	OAuth    OAuthConfig
}

type ServerConfig struct {
	Address string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	Issuer             string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() (*Config, error) {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// PostgreSQL port'unu int'e çevir
	dbPort, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))

	// Redis DB numarasını int'e çevir
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	// SMTP port'unu int'e çevir
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	// JWT TTL'leri parse et
	accessTTL, _ := time.ParseDuration(os.Getenv("JWT_ACCESS_TTL"))
	refreshTTL, _ := time.ParseDuration(os.Getenv("JWT_REFRESH_TTL"))

	return &Config{
		Server: ServerConfig{
			Address: ":8080",
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     dbPort,
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Address:  os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			AccessTokenSecret:  os.Getenv("JWT_ACCESS_SECRET"),
			RefreshTokenSecret: os.Getenv("JWT_REFRESH_SECRET"),
			AccessTokenTTL:     accessTTL,
			RefreshTokenTTL:    refreshTTL,
			Issuer:             os.Getenv("JWT_ISSUER"),
		},
		SMTP: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     smtpPort,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		},
	}, nil
}
