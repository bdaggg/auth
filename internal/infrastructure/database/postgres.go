package database

import (
	"fmt"

	"auth-service/internal/config"
	"auth-service/internal/domain/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresConnection(cfg config.DatabaseConfig) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres bağlantısı açılamadı: %v", err)
	}

	// Otomatik migrasyon
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Subscription{},
	)
	if err != nil {
		return nil, fmt.Errorf("migrasyon hatası: %v", err)
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the underlying database connection
func (p *PostgresDB) GetDB() *gorm.DB {
	return p.db
}
