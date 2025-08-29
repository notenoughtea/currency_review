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
	data *dto.CurrencyRates
}

func GetRatesService() RatesService {
	dbInst, err := db.Connect()
	if err != nil {
		log.Println("Ошибка подключения к базе")
	}
	var rate dto.CurrencyRates
	dbInst.Order("created_at desc").First(&rate)

	return &ratesService{
		data: &dto.CurrencyRates{
			ID:                 rate.ID,
			Result:             rate.Result,
			Documentation:      rate.Documentation,
			TermsOfUse:         rate.TermsOfUse,
			TimeLastUpdateUnix: rate.TimeLastUpdateUnix,
			TimeLastUpdateUtc:  rate.TimeLastUpdateUtc,
			TimeNextUpdateUnix: rate.TimeNextUpdateUnix,
			TimeNextUpdateUtc:  rate.TimeNextUpdateUtc,
			BaseCode:           rate.BaseCode,
			ConversionRates:    rate.ConversionRates,
			CreatedAt:          rate.CreatedAt,
			UpdatedAt:          rate.UpdatedAt,
		},
	}
}

func (s *ratesService) GetAll() *dto.CurrencyRates {
	return s.data
}

func (s *ratesService) Get(code string) (float64, bool) {
	return s.data.GetRate(code)
}
