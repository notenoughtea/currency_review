package dto

import "time"

type Rates map[string]float64

type CurrencyRates struct {
	ID                 uint   `gorm:"primaryKey"`
	Result             string `json:"result"`
	Documentation      string `json:"documentation"`
	TermsOfUse         string `json:"terms_of_use"`
	TimeLastUpdateUnix int64  `json:"time_last_update_unix"`
	TimeLastUpdateUtc  string `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64  `json:"time_next_update_unix"`
	TimeNextUpdateUtc  string `json:"time_next_update_utc"`
	BaseCode           string `json:"base_code"`
	ConversionRates    Rates  `json:"conversion_rates" gorm:"type:jsonb;serializer:json"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (c *CurrencyRates) GetRate(code string) (float64, bool) {
	v, ok := c.ConversionRates[code]
	return v, ok
}
