package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/notenoughtea/currency_review/gateway/internal/config"
	"github.com/notenoughtea/currency_review/gateway/internal/handler"
	"github.com/notenoughtea/currency_review/gateway/internal/logger"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func observeMetrics(counter *prometheus.CounterVec, hist *prometheus.HistogramVec, r *http.Request, status int, start time.Time) {
	path := r.URL.Path
	labels := prometheus.Labels{"method": r.Method, "path": path, "status": fmt.Sprintf("%d", status)}
	counter.With(labels).Inc()
	hist.With(labels).Observe(time.Since(start).Seconds())
}

func main() {
	// включаем логи
	logger.Init()

	// подключаем prometheus
	config.Load()
	addr := fmt.Sprintf("%s:%d", config.GetServerConfig().Host, config.GetServerConfig().Port)

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gateway",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "gateway",
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request duration seconds",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)
	register := func(c prometheus.Collector) {
		if err := prometheus.DefaultRegisterer.Register(c); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
				return
			}
		}
	}
	register(requestCounter)
	register(requestDuration)
	register(collectors.NewBuildInfoCollector())
	register(collectors.NewGoCollector())
	register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/code/", handler.GetByCodeHandler)
	mux.Handle("/metrics", promhttp.Handler())

	wrap := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// capture status code
			ww := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			defer func() {
				if rec := recover(); rec != nil {
					http.Error(ww, "internal error", http.StatusInternalServerError)
					observeMetrics(requestCounter, requestDuration, r, ww.status, start)
				}
			}()
			h.ServeHTTP(ww, r)
			observeMetrics(requestCounter, requestDuration, r, ww.status, start)
		})
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      wrap(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Log.Infof("listen on %s", addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Log.Errorf("graceful shutdown failed: %v", err)
			_ = srv.Close()
		} else {
			logger.Log.Infof("server stopped gracefully")
		}
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("server error:", err)
		}
	}

	logger.Log.Info("server exited")

}
