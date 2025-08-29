package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/worker"
)

// "github.com/notenoughtea/currency_review/currency/internal/worker"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.GetRatesDaily(ctx, 12*time.Hour)
	}()

	wg.Wait()
}
