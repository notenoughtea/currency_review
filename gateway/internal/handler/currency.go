package handler

import (
	"net/http"
	"time"

	"github.com/notenoughtea/currency_review/gateway/internal/dto"
	"github.com/notenoughtea/currency_review/gateway/internal/logger"
)

func (s *controller) GetCurrencyRates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, err := s.currencyService.GetCurrencyRates(r.Context())
	if err != nil {
		logger.Log.Errorf("get rates error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *controller) GetRatesByDates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dateFromStr := r.URL.Query().Get("date_from")
	dateToStr := r.URL.Query().Get("date_to")
	if dateFromStr == "" || dateToStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing required params"})
		return
	}

	dateFrom, err := time.Parse("2006-01-02", dateFromStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date_from format"})
		return
	}
	dateTo, err := time.Parse("2006-01-02", dateToStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date_to format"})
		return
	}
	if dateFrom.After(dateTo) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "date_from must be <= date_to"})
		return
	}

	req := &dto.ParsedCurrencyRequest{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	data, err := s.currencyService.GetRatesByDates(req)
	if err != nil {
		logger.Log.Errorf("get rates by dates error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, data)
}
