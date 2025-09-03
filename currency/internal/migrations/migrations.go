package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"gorm.io/gorm"
)

func MigrateCurrencyTable(db *gorm.DB) error {

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "2023081501",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&dto.CurrencyRates{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("currency_rates")
			},
		},
	})
	if err := m.Migrate(); err != nil {
		return err
	}

	logger.Log.Errorf("Миграции успешно выполнены")
	return nil
}
