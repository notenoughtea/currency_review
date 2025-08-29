package currclnt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
)

func GetRates() dto.CurrencyRates {
	token := config.GetToken()
	requestString := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", token)
	resp, err := http.Get(requestString)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var cr dto.CurrencyRates
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		panic(err)
	}
	log.Println("Получены свежие котировки")
	return cr
}
