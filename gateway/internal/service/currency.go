package service

import (
	"context"

	clients "github.com/notenoughtea/currency_review/gateway/internal/clients/gRPC"
	"github.com/notenoughtea/currency_review/gateway/internal/dto"
)

type CurrencyServiceInterface interface {
	GetCurrencyRates(ctx context.Context) (any, error)
	GetRatesByDates(req *dto.ParsedCurrencyRequest) ([]clients.CurrencyRate, error)
}

type CurrencyService struct{}

func NewCurrencyService() *CurrencyService {
	return &CurrencyService{}
}

func (s *CurrencyService) GetCurrencyRates(ctx context.Context) (any, error) {
	rates := clients.GetAllRatesHandler()
	return rates, nil
}

func (s *CurrencyService) GetRatesByDates(req *dto.ParsedCurrencyRequest) ([]clients.CurrencyRate, error) {
	return clients.GetRatesByDates(req)
}
