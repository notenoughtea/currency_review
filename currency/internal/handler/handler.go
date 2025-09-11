package handler

import (
	"context"

	"github.com/notenoughtea/currency_review/currency/internal/service"
	"github.com/notenoughtea/currency_review/pkg"
)

type RatesHandler struct {
	pkg.UnimplementedRatesServiceServer
	svc service.RatesService
}

func NewRatesHandler(svc service.RatesService) *RatesHandler {
	return &RatesHandler{svc: svc}
}

func (h *RatesHandler) GetAllRates(ctx context.Context, _ *pkg.Empty) (*pkg.CurrencyRates, error) {
	cr := h.svc.GetAll()

	rates := make([]*pkg.Rate, 0, len(cr.Rates))
	for _, r := range cr.Rates {
		rates = append(rates, &pkg.Rate{
			Date: r.Date.Format("2006-01-02"),
			Rate: r.Rate,
		})
	}

	return &pkg.CurrencyRates{
		Currency: cr.Currency,
		Rates:    rates,
	}, nil
}
