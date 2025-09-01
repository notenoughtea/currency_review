package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/worker"
)

// "github.com/notenoughtea/currency_review/currency/internal/worker"

func main() {
	config.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Запуск сервиса по ежедневному запросу курсов валют")
		worker.GetRatesDaily(ctx, 12*time.Hour)
	}()

	wg.Wait()
}
