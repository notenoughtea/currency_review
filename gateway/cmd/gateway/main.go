package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/notenoughtea/currency_review/gateway/internal/config"
	"github.com/notenoughtea/currency_review/gateway/internal/handler"
	"github.com/notenoughtea/currency_review/gateway/internal/logger"
)

func main() {
	// включаем логи
	logger.Init()

	config.Load()
	addr := fmt.Sprintf("%s:%d", config.GetServerConfig().Host, config.GetServerConfig().Port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/code/", handler.GetByCodeHandler)

	wrap := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
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
