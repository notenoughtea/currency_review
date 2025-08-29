package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/notenoughtea/currency_review/gateway/internal/clients"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	rates := clients.GetAllRatesHandler()
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.MarshalIndent(rates, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func GetByCodeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "missing code in path", http.StatusBadRequest)
		return
	}
	code := strings.ToUpper(strings.TrimSpace(parts[2]))
	rate, err := clients.GetRateHandler(code)
	if err != nil {
		fmt.Fprintf(w, "не найдено котировок с кодом: %s\n", code)
		log.Println(err)
	}
	if rate != 0 {
		fmt.Fprintf(w, "Вы запросили курс EUR к %v\n", code)
		fmt.Fprintf(w, "Ваш курс 1 к %v\n", rate)
		return
	}
	fmt.Fprintf(w, "Вы запросили курс EUR к %v\n", code)
	fmt.Fprintf(w, "Ваш курс не найден %v\n", rate)
}
