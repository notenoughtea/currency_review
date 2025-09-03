package repository

import (
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"gorm.io/gorm"
)

type RatesRepository interface {
	StoreRates(newRates dto.CurrencyRates) error
	GetLatestRates() (*dto.CurrencyRates, error)
}

type ratesRepository struct {
	db *gorm.DB
}

func NewRatesRepository(db *gorm.DB) RatesRepository {
	return &ratesRepository{db: db}
}

func (r *ratesRepository) Migrate(newRates dto.CurrencyRates) error {
	r.db.AutoMigrate(&dto.CurrencyRates{})
	logger.Log.Infof("Созданы записи в ДБ: %v, %v %v", newRates.BaseCode, len(newRates.ConversionRates), "курсов")
	return nil
}

func (r *ratesRepository) StoreRates(newRates dto.CurrencyRates) error {
	r.db.AutoMigrate(&dto.CurrencyRates{})
	if err := r.db.Create(&newRates).Error; err != nil {
		return err
	}
	logger.Log.Infof("Созданы записи в ДБ: %v, %v %v", newRates.BaseCode, len(newRates.ConversionRates), "курсов")
	return nil
}

func (r *ratesRepository) GetLatestRates() (*dto.CurrencyRates, error) {
	var rate dto.CurrencyRates
	if err := r.db.Order("created_at desc").First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}
