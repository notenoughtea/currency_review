package repository

import (
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"gorm.io/gorm"
)

type RatesRepository interface {
	StoreRates(newRates dto.CurrencyRates) error
	GetLatestRates() (*dto.CurrencyRates, error)
	GetRatesByDates(*dto.CurrencyRequest) (*[]dto.CurrencyRates, error)
}

type ratesRepository struct {
	db *gorm.DB
}

func NewRatesRepository(db *gorm.DB) RatesRepository {
	return &ratesRepository{db: db}
}

type CurrencyRatesDB struct {
	ID        uint     `gorm:"primaryKey"`
	Currency  string   `gorm:"uniqueIndex;not null"`
	Rates     []RateDB `gorm:"foreignKey:CurrencyRatesID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RateDB struct {
	ID              uint      `gorm:"primaryKey"`
	CurrencyRatesID uint      `gorm:"index;not null"`
	Date            time.Time `gorm:"index;not null"`
	Rate            float32   `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func toDBModel(in dto.CurrencyRates) CurrencyRatesDB {
	out := CurrencyRatesDB{Currency: in.Currency}
	out.Rates = make([]RateDB, 0, len(in.Rates))
	for _, r := range in.Rates {
		out.Rates = append(out.Rates, RateDB{Date: r.Date, Rate: r.Rate})
	}
	return out
}

func toDTOModel(in CurrencyRatesDB) dto.CurrencyRates {
	out := dto.CurrencyRates{Currency: in.Currency}
	out.Rates = make([]dto.Rate, 0, len(in.Rates))
	for _, r := range in.Rates {
		out.Rates = append(out.Rates, dto.Rate{Date: r.Date, Rate: r.Rate})
	}
	return out
}

func (r *ratesRepository) Migrate(newRates dto.CurrencyRates) error {
	baseCode := config.GetDefaultBaseCurrency()
	if err := r.db.AutoMigrate(&CurrencyRatesDB{}, &RateDB{}); err != nil {
		return err
	}
	logger.Log.Infof("Миграции применены: %v", baseCode)
	return nil
}

func (r *ratesRepository) StoreRates(newRates dto.CurrencyRates) error {
	baseCode := config.GetDefaultBaseCurrency()
	if err := r.db.AutoMigrate(&CurrencyRatesDB{}, &RateDB{}); err != nil {
		return err
	}
	dbModel := toDBModel(newRates)
	if err := r.db.Create(&dbModel).Error; err != nil {
		return err
	}
	logger.Log.Infof("Созданы записи в ДБ: %v, %v %v", baseCode, len(newRates.Rates), "курсов")
	return nil
}

func (r *ratesRepository) GetLatestRates() (*dto.CurrencyRates, error) {
	var rec CurrencyRatesDB
	if err := r.db.Preload("Rates").Order("created_at desc").First(&rec).Error; err != nil {
		return nil, err
	}
	dtoRec := toDTOModel(rec)
	return &dtoRec, nil
}

func (r *ratesRepository) GetRatesByDates(req *dto.CurrencyRequest) (*[]dto.CurrencyRates, error) {
	var rows []CurrencyRatesDB
	if err := r.db.Preload("Rates").Where("created_at BETWEEN ? AND ?", req.DateFrom, req.DateTo).
		Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]dto.CurrencyRates, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDTOModel(row))
	}
	return &out, nil
}
