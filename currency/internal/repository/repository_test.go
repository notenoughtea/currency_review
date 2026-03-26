package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/notenoughtea/currency_review/currency/internal/dto"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&CurrencyRatesDB{}, &RateDB{}))
	return db
}

func TestStoreAndGetLatestRates(t *testing.T) {
	db := setupTestDB(t)
	r := NewRatesRepository(db)

	rates1 := dto.CurrencyRates{Currency: "USD", Rates: []dto.Rate{{Date: time.Now(), Rate: 0.9}}}
	require.NoError(t, r.StoreRates(rates1))

	got1, err := r.GetLatestRates()
	require.NoError(t, err)
	require.Equal(t, "USD", got1.Currency)
	require.Len(t, got1.Rates, 1)

	rates2 := dto.CurrencyRates{Currency: "USD", Rates: []dto.Rate{{Date: time.Now(), Rate: 0.95}}}
	require.NoError(t, r.StoreRates(rates2))

	got2, err := r.GetLatestRates()
	require.NoError(t, err)
	require.Equal(t, float32(0.95), got2.Rates[0].Rate)
}
