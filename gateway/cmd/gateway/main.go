package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/notenoughtea/currency_review/gateway/internal/config"
	"github.com/notenoughtea/currency_review/gateway/internal/handler"
)

func main() {
	// r, err := clients.GetRateHandler("EUR")
	// if err != nil {
	// 	log.Println("котировка не найдена")
	// }
	// fmt.Println(r)

	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/code/", handler.GetByCodeHandler)

	conf := config.GetServerConfig()
	portString := fmt.Sprintf(":%v", conf.Port)
	log.Printf("Starting server at port %v", portString)
	err := http.ListenAndServe(portString, nil)
	if err != nil {
		log.Println("Error starting the server:", err)
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
