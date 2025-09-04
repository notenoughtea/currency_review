package currclnt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	currClientRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "currency",
			Subsystem: "external_client",
			Name:      "requests_total",
			Help:      "Total requests to external currency API",
		},
		[]string{"endpoint", "code"},
	)
	currClientDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "currency",
			Subsystem: "external_client",
			Name:      "request_duration_seconds",
			Help:      "Duration of requests to external currency API",
			Buckets:   []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"endpoint"},
	)
)

func init() {
	prometheus.MustRegister(currClientRequests, currClientDuration)
}

func GetRates() dto.CurrencyRates {
	token := config.GetToken()
	requestString := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", token)
	start := time.Now()
	resp, err := http.Get(requestString)
	if err != nil {
		currClientRequests.WithLabelValues("latest_USD", "error").Inc()
		currClientDuration.WithLabelValues("latest_USD").Observe(time.Since(start).Seconds())
		logger.Log.Fatal(err)
	}
	defer resp.Body.Close()
	currClientRequests.WithLabelValues("latest_USD", fmt.Sprintf("%d", resp.StatusCode)).Inc()
	currClientDuration.WithLabelValues("latest_USD").Observe(time.Since(start).Seconds())

	var cr dto.CurrencyRates
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		logger.Log.Fatal(err)
	}
	logger.Log.Info("Получены свежие котировки")
	return cr
}
