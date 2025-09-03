package db

import (
	"fmt"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	cfg := config.GetDbConfig()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s connect_timeout=10",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode, cfg.TimeZone,
	)

	var db *gorm.DB
	var err error

	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, e := db.DB()
			if e == nil && sqlDB.Ping() == nil {
				logger.Log.Info("Подключение к БД успешно")
				return db, nil
			}
			if e != nil {
				err = e
			}
		}
		logger.Log.Errorf("подключение к БД неуспешно, попытка %d: %v", i, err)
		time.Sleep(time.Duration(i) * time.Second)
	}

	return nil, err
}
