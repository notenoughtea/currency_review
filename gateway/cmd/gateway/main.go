package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/notenoughtea/currency_review/gateway/internal/config"
	"github.com/notenoughtea/currency_review/gateway/internal/handler"
)

func main() {
	// r, err := clients.GetRateHandler("EUR")
	// if err != nil {
	// 	log.Println("котировка не найдена")
	// }
	// fmt.Pcon
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

	log.Printf("listen on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	// тут про логи
	// file, err := os.OpenFile("../../../logs/logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	// logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
	// 	Level: slog.LevelInfo,
	// }))
	// slog.SetDefault(logger)
	// slog.Info("Gateway: запущен", "version", "1.0.0")
	// slog.Warn("Gateway: Предупреждение", "disk", "80%")
	// slog.Error("Gateway: Ошибка", "code", 500)
}
