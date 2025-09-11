package main

import (
	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/db"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"github.com/notenoughtea/currency_review/currency/internal/migrations"
)

func main() {
	config.Load()
	conn, err := db.Connect()
	if err != nil {
		logger.Log.Fatalf("Ошибка подключения к базе: %v", err)
	}

	if err := migrations.MigrateCurrencyTable(conn); err != nil {
		logger.Log.Fatalf("Ошибка миграции: %v", err)
	}
}
