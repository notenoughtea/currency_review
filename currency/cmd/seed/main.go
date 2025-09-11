package main

import (
	"math/rand"
	"sort"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/db"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"github.com/notenoughtea/currency_review/currency/internal/migrations"
	"github.com/notenoughtea/currency_review/currency/internal/repository"
)

func randDates(start, end time.Time, n int) []time.Time {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	days := int(end.Sub(start).Hours() / 24)
	if days < 1 || n <= 0 {
		return nil
	}
	set := make(map[string]struct{}, n)
	out := make([]time.Time, 0, n)
	for len(out) < n {
		offset := rng.Intn(days + 1)
		d := start.AddDate(0, 0, offset).Truncate(24 * time.Hour).UTC()
		k := d.Format("2006-01-02")
		if _, ok := set[k]; ok {
			continue
		}
		set[k] = struct{}{}
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Before(out[j]) })
	return out
}

func main() {
	config.Load()
	conn, err := db.Connect()
	if err != nil {
		logger.Log.Fatalf("Ошибка подключения к базе: %v", err)
	}

	if err := migrations.MigrateCurrencyTable(conn); err != nil {
		logger.Log.Fatalf("Ошибка миграции: %v", err)
	}

	repo := repository.NewRatesRepository(conn)

	now := time.Now().UTC()
	start := now.AddDate(0, 0, -120)
	dates := randDates(start, now, 90)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type cdef struct {
		Code string
		Base float64
		Vol  float64
	}

	defs := []cdef{
		{Code: "USD", Base: 90.0, Vol: 0.35},
		{Code: "EUR", Base: 95.0, Vol: 0.30},
		{Code: "GBP", Base: 110.0, Vol: 0.40},
	}

	seed := make([]dto.CurrencyRates, 0, len(defs))
	for _, d := range defs {
		cr := dto.CurrencyRates{Currency: d.Code}
		cr.Rates = make([]dto.Rate, 0, len(dates))
		level := d.Base
		for _, day := range dates {
			level += rng.NormFloat64() * d.Vol
			if level < 0 {
				level = d.Base
			}
			cr.Rates = append(cr.Rates, dto.Rate{
				Date: day,
				Rate: float32(level),
			})
		}
		seed = append(seed, cr)
	}

	for _, r := range seed {
		if err := repo.StoreRates(r); err != nil {
			logger.Log.Fatalf("Ошибка наполнения БД: %v", err)
		}
	}

	logger.Log.Infof("Сидер: создано %d записей валют, по %d дат на каждую", len(seed), len(dates))
}
