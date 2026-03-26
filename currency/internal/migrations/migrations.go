package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"gorm.io/gorm"
)

func MigrateCurrencyTable(db *gorm.DB) error {

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "2023081501",
			Migrate: func(tx *gorm.DB) error {
				type Rate struct {
					ID              uint      `gorm:"primaryKey"`
					CurrencyRatesID uint      `gorm:"index;not null"`
					Date            time.Time `gorm:"index;not null"`
					Rate            float32   `gorm:"not null"`
					CreatedAt       time.Time
					UpdatedAt       time.Time
				}

				type CurrencyRates struct {
					ID        uint   `gorm:"primaryKey"`
					Currency  string `gorm:"uniqueIndex;not null"`
					Rates     []Rate `gorm:"foreignKey:CurrencyRatesID;constraint:OnDelete:CASCADE"`
					CreatedAt time.Time
					UpdatedAt time.Time
				}

				return tx.AutoMigrate(&CurrencyRates{}, &Rate{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("rates", "currency_rates")
			},
		},
	})
	if err := m.Migrate(); err != nil {
		return err
	}

	logger.Log.Info("Миграции успешно выполнены")
	return nil
}
