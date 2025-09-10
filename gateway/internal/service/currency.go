package service

import (
	"context"

	clients "github.com/notenoughtea/currency_review/gateway/internal/clients/gRPC"
)

type CurrencyServiceInterface interface {
	GetCurrencyRates(ctx context.Context) (any, error)
}

type CurrencyService struct{}

func NewCurrencyService() *CurrencyService {
	return &CurrencyService{}
}

func (s *CurrencyService) GetCurrencyRates(ctx context.Context) (any, error) {
	rates := clients.GetAllRatesHandler()
	return rates, nil
}
