package repository

package repository

import (
	"fmt"

	"github.com/notenoughtea/currency_review/currency/internal/dto"
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

func (r *ratesRepository) StoreRates(newRates dto.CurrencyRates) error {
	r.db.AutoMigrate(&dto.CurrencyRates{})
	if err := r.db.Create(&newRates).Error; err != nil {
		return err
	}
	fmt.Println("Созданы записи в ДБ:", newRates.BaseCode, len(newRates.ConversionRates), "курсов")
	return nil
}

func (r *ratesRepository) GetLatestRates() (*dto.CurrencyRates, error) {
	var rate dto.CurrencyRates
	if err := r.db.Order("created_at desc").First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}
