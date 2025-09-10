package handler

import (
	"encoding/json"
	"net/http"

	"github.com/notenoughtea/currency_review/gateway/internal/service"
	"github.com/sirupsen/logrus"
)

type controller struct {
	authService     service.AuthServiceInterface
	currencyService service.CurrencyServiceInterface
	logger          *logrus.Logger
}

func RegisterRoutes(authSvc service.AuthServiceInterface, currencySvc service.CurrencyServiceInterface, mux *http.ServeMux, log *logrus.Logger) controller {
	cntrl := controller{authService: authSvc, currencyService: currencySvc, logger: log}
	mux.HandleFunc("/ping", cntrl.ping)
	mux.HandleFunc("/api/v1/rate", cntrl.GetCurrencyRates)
	mux.HandleFunc("/api/v1/login", cntrl.Login)
	mux.HandleFunc("/api/v1/register", cntrl.Register)
	mux.HandleFunc("/api/v1/logout", cntrl.Logout)
	return cntrl
}

func (s *controller) ping(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "pong"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
