package service

import (
	"fmt"
	"log"

	"github.com/notenoughtea/currency_review/currency/internal/clients/currclnt"
	"github.com/notenoughtea/currency_review/currency/internal/db"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/migrations"
	"gorm.io/gorm"
)

func StoreRates(db *gorm.DB, newRates dto.CurrencyRates) {

	db.AutoMigrate(&dto.CurrencyRates{})
	if err := db.Create(&newRates).Error; err != nil {
		panic(err)
	}
	fmt.Println("Созданы записи в ДБ:", newRates.BaseCode, len(newRates.ConversionRates), "курсов")
}

func HandleRates() {
	dbInst, err := db.Connect()
	if err != nil {
		log.Println("Ошибка подключения к базе")
	}
	err = migrations.MigrateCurrencyTable(dbInst)
	if err != nil {
		log.Println("Миграция не прошла")
	}
	newRates := currclnt.GetRates()
	StoreRates(dbInst, newRates)
}

type RatesService interface {
	GetAll() *dto.CurrencyRates
	Get(code string) (float64, bool)
}

type ratesService struct {
	db *gorm.DB
}

func GetRatesService() RatesService {
	dbInst, err := db.Connect()
	if err != nil {
		log.Println("Ошибка подключения к базе")
	}
	return &ratesService{db: dbInst}
}

func (s *ratesService) GetAll() *dto.CurrencyRates {
	var rate dto.CurrencyRates
	s.db.Order("created_at desc").First(&rate)
	return &rate
}

func (s *ratesService) Get(code string) (float64, bool) {
	var rate dto.CurrencyRates
	s.db.Order("created_at desc").First(&rate)
	return rate.GetRate(code)
}
