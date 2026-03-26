package dto

import "time"

type ParsedCurrencyRequest struct {
	DateFrom time.Time
	DateTo   time.Time
}

type CurrencyResponse struct {
	Currency string
	Rates    []CurrencyRate
}

type CurrencyRate struct {
	Rate float32
	Date time.Time
}
