package worker

import (
	"context"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"github.com/notenoughtea/currency_review/currency/internal/service"
)

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
		}
	}
}
