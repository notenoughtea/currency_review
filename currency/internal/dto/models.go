package dto

import (
	"time"
)

type CurrencyRequest struct {
	BaseCurrency   string
	TargetCurrency string
	DateFrom       time.Time
	DateTo         time.Time
}

type CurrencyRates struct {
	Currency string
	Rates    []Rate
}

type Rate struct {
	Date time.Time
	Rate float32
}
