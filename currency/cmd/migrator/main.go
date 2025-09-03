package main

import (
	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/db"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
)

func main() {
	config.Load()
	conn, err := db.Connect()
	if err != nil {
		logger.Log.Fatalf("Ошибка подключения к базе: ", err)
	}

	if err := conn.AutoMigrate(
		&dto.CurrencyRates{},
	); err != nil {
		logger.Log.Fatalf("Ошибка миграции:", err)
	}

	logger.Log.Info("Миграции успешно выполнены")
}
