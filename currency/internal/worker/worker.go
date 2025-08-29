package worker

import (
	"context"
	"log"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/service"
)

func GetRatesDaily(ctx context.Context, period time.Duration) {
	service.HandleRates()
	t := time.NewTicker(period)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Printf("Отмена периодического запроса курсов")
			return
		case <-t.C:
			service.HandleRates()
			log.Printf("Обновлен курс, %v", time.Now())
		}
	}
}
