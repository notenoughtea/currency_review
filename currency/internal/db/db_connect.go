package db

import (
	"fmt"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	cfg := config.GetDbConfig()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode, cfg.TimeZone,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
