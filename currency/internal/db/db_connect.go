package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var registerDbMetricsOnce sync.Once

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
				stats := func() *gorm.DB {
					return db
				}
				registerDbMetricsOnce.Do(func() {
					prometheus.MustRegister(
						prometheus.NewGaugeFunc(prometheus.GaugeOpts{Namespace: "currency", Subsystem: "db_pool", Name: "open_connections", Help: "DB open connections"}, func() float64 {
							s, _ := stats().DB()
							return float64(s.Stats().OpenConnections)
						}),
						prometheus.NewGaugeFunc(prometheus.GaugeOpts{Namespace: "currency", Subsystem: "db_pool", Name: "in_use", Help: "DB connections in use"}, func() float64 {
							s, _ := stats().DB()
							return float64(s.Stats().InUse)
						}),
						prometheus.NewGaugeFunc(prometheus.GaugeOpts{Namespace: "currency", Subsystem: "db_pool", Name: "idle", Help: "DB idle connections"}, func() float64 {
							s, _ := stats().DB()
							return float64(s.Stats().Idle)
						}),
						prometheus.NewGaugeFunc(prometheus.GaugeOpts{Namespace: "currency", Subsystem: "db_pool", Name: "wait_count", Help: "DB wait count"}, func() float64 {
							s, _ := stats().DB()
							return float64(s.Stats().WaitCount)
						}),
						prometheus.NewGaugeFunc(prometheus.GaugeOpts{Namespace: "currency", Subsystem: "db_pool", Name: "wait_duration_seconds", Help: "DB wait duration seconds"}, func() float64 {
							s, _ := stats().DB()
							return s.Stats().WaitDuration.Seconds()
						}),
					)
				})
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
