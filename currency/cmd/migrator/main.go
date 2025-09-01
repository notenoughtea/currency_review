package main

import (
	"log"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/db"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
)

func main() {
	config.Load()
	conn, err := db.Connect()
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}

	if err := conn.AutoMigrate(
		&dto.CurrencyRates{},
	); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	log.Println("Миграции успешно выполнены")
}
