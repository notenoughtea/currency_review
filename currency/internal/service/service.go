package service

import (
	"log"

	"github.com/notenoughtea/currency_review/currency/internal/clients/currclnt"
	"github.com/notenoughtea/currency_review/currency/internal/db"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/migrations"
	"github.com/notenoughtea/currency_review/currency/internal/repository"
	"gorm.io/gorm"
)

type RatesService interface {
	GetAll() *dto.CurrencyRates
	Get(code string) (float64, bool)
	HandleRates()
}

type ratesService struct {
	db   *gorm.DB
	repo repository.RatesRepository
}

func GetRatesService() RatesService {
	dbInst, err := db.Connect()
	if err != nil {
		log.Println("Ошибка подключения к базе")
	}
	return &ratesService{
		db:   dbInst,
		repo: repository.NewRatesRepository(dbInst),
	}
}

func (s *ratesService) HandleRates() {
	if err := migrations.MigrateCurrencyTable(s.db); err != nil {
		log.Println("Миграция не прошла")
		return
	}
	newRates := currclnt.GetRates()
	if err := s.repo.StoreRates(newRates); err != nil {
		log.Println("Ошибка при получении курсов:", err)
	}
}

func (s *ratesService) GetAll() *dto.CurrencyRates {
	rate, err := s.repo.GetLatestRates()
	if err != nil {
		log.Println("Ошибка при получении курсов:", err)
		return nil
	}
	return rate
}

func (s *ratesService) Get(code string) (float64, bool) {
	rate, err := s.repo.GetLatestRates()
	if err != nil {
		log.Println("Ошибка при получении курсов:", err)
		return 0, false
	}
	return rate.GetRate(code)
}
