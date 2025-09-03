package currclnt

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
)

func GetRates() dto.CurrencyRates {
	token := config.GetToken()
	requestString := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", token)
	resp, err := http.Get(requestString)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer resp.Body.Close()

	var cr dto.CurrencyRates
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		logger.Log.Fatal(err)
	}
	logger.Log.Info("Получены свежие котировки")
	return cr
}
