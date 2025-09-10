package handler

import (
	"net/http"

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
