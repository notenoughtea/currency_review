package worker

import (
	"context"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"github.com/notenoughtea/currency_review/currency/internal/service"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	workerRuns = prometheus.NewCounterVec(
		prometheus.CounterOpts{Namespace: "currency", Subsystem: "worker", Name: "jobs_processed_total", Help: "Processed jobs"},
		[]string{"job", "result"},
	)
	workerLastRun = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Namespace: "currency", Subsystem: "worker", Name: "last_run_timestamp_seconds", Help: "Last run UNIX timestamp"},
		[]string{"job"},
	)
)

func init() {
	prometheus.MustRegister(workerRuns, workerLastRun)
}

func GetRatesDaily(ctx context.Context, period time.Duration) {
	service.GetRatesService().HandleRates()
	t := time.NewTicker(period)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Отмена периодического запроса курсов")
			return
		case <-t.C:
			service.GetRatesService().HandleRates()
			logger.Log.Infof("Обновлен курс, %v", time.Now())
			workerRuns.WithLabelValues("rates_update", "success").Inc()
			workerLastRun.WithLabelValues("rates_update").SetToCurrentTime()
		}
	}
}
