package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/notenoughtea/currency_review/currency/internal/dto"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&dto.CurrencyRates{}))
	return db
}

func TestStoreAndGetLatestRates(t *testing.T) {
	db := setupTestDB(t)
	r := NewRatesRepository(db)

	rates1 := dto.CurrencyRates{BaseCode: "USD", ConversionRates: dto.Rates{"EUR": 0.9}}
	require.NoError(t, r.StoreRates(rates1))

	got1, err := r.GetLatestRates()
	require.NoError(t, err)
	require.Equal(t, "USD", got1.BaseCode)
	require.InDelta(t, 0.9, got1.ConversionRates["EUR"], 1e-9)

	rates2 := dto.CurrencyRates{BaseCode: "USD", ConversionRates: dto.Rates{"EUR": 0.95}}
	require.NoError(t, r.StoreRates(rates2))

	got2, err := r.GetLatestRates()
	require.NoError(t, err)
	require.InDelta(t, 0.95, got2.ConversionRates["EUR"], 1e-9)
}
