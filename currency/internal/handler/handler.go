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

func (h *RatesHandler) GetRate(ctx context.Context, req *pkg.GetRateRequest) (*pkg.GetRateResponse, error) {
	v, ok := h.svc.Get(req.GetCode())
	return &pkg.GetRateResponse{Value: v, Found: ok}, nil
}

func (h *RatesHandler) GetAllRates(ctx context.Context, _ *pkg.Empty) (*pkg.CurrencyRates, error) {
	cr := h.svc.GetAll()
	return &pkg.CurrencyRates{
		Id:                 uint64(cr.ID),
		Result:             cr.Result,
		Documentation:      cr.Documentation,
		TermsOfUse:         cr.TermsOfUse,
		TimeLastUpdateUnix: cr.TimeLastUpdateUnix,
		TimeLastUpdateUtc:  cr.TimeLastUpdateUtc,
		TimeNextUpdateUnix: cr.TimeNextUpdateUnix,
		TimeNextUpdateUtc:  cr.TimeNextUpdateUtc,
		BaseCode:           cr.BaseCode,
		ConversionRates:    map[string]float64(cr.ConversionRates),
		CreatedAt:          cr.CreatedAt.Unix(),
		UpdatedAt:          cr.UpdatedAt.Unix(),
	}, nil
}
