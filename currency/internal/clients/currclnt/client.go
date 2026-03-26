package currclnt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	baseCode := strings.ToLower(config.GetDefaultBaseCurrency())
	baseURL := strings.TrimRight(config.GetOuterUrl(), "/")
	url := baseURL
	if !strings.HasSuffix(baseURL, ".json") {
		url = fmt.Sprintf("%s/%s.json", baseURL, baseCode)
	}
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		currClientRequests.WithLabelValues(baseCode, "error").Inc()
		currClientDuration.WithLabelValues(baseCode).Observe(time.Since(start).Seconds())
		logger.Log.Fatal(err)
	}
	defer resp.Body.Close()
	currClientRequests.WithLabelValues(baseCode, fmt.Sprintf("%d", resp.StatusCode)).Inc()
	currClientDuration.WithLabelValues(baseCode).Observe(time.Since(start).Seconds())

	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()

	var generic map[string]interface{}
	if err := dec.Decode(&generic); err != nil {
		logger.Log.Fatal(err)
	}

	actualBase := baseCode
	innerAny, ok := generic[baseCode]
	if !ok {
		for k, v := range generic {
			kl := strings.ToLower(k)
			if kl == "date" {
				continue
			}
			actualBase = kl
			innerAny = v
			break
		}
	}

	innerMap, ok := innerAny.(map[string]interface{})
	if !ok {
		logger.Log.Fatalf("неожиданный формат JSON для базовой валюты %s", actualBase)
	}

	useDate := time.Now().UTC()
	if dv, ok := innerMap["date"]; ok {
		if ds, ok := dv.(string); ok {
			if t, err := time.Parse(time.RFC3339, ds); err == nil {
				useDate = t.UTC()
			} else if t2, err2 := time.Parse("2006-01-02", ds); err2 == nil {
				useDate = t2.UTC()
			}
		}
	}

	rates := make([]dto.Rate, 0, len(innerMap))
	for k, rawVal := range innerMap {
		if strings.ToLower(k) == "date" {
			continue
		}
		var f64 float64
		switch v := rawVal.(type) {
		case json.Number:
			ff, err := v.Float64()
			if err != nil {
				continue
			}
			f64 = ff
		case float64:
			f64 = v
		case string:
			ff, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			f64 = ff
		default:
			continue
		}
		rates = append(rates, dto.Rate{Date: useDate, Rate: float32(f64)})
	}

	cr := dto.CurrencyRates{
		Currency: strings.ToUpper(actualBase),
		Rates:    rates,
	}
	logger.Log.Info("Получены свежие котировки")
	return cr
}
