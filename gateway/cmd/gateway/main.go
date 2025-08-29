package main

import (
	"log/slog"
	"os"

	"github.com/notenoughtea/currency_review/gateway/internal/clients"
)

func main() {
	clients.GPRCclnt()

	file, err := os.OpenFile("../../../logs/logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Gateway: запущен", "version", "1.0.0")
	slog.Warn("Gateway: Предупреждение", "disk", "80%")
	slog.Error("Gateway: Ошибка", "code", 500)
}
