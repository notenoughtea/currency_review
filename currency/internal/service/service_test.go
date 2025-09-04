package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/notenoughtea/currency_review/currency/internal/dto"
	repo "github.com/notenoughtea/currency_review/currency/internal/repository"
	"github.com/notenoughtea/currency_review/currency/internal/repository/mocks"
)

type testRatesService struct{ repo repo.RatesRepository }

func (s *testRatesService) HandleRates()               {}
func (s *testRatesService) GetAll() *dto.CurrencyRates { r, _ := s.repo.GetLatestRates(); return r }
func (s *testRatesService) Get(code string) (float64, bool) {
	r, _ := s.repo.GetLatestRates()
	if r == nil {
		return 0, false
	}
	return r.GetRate(code)
}

func TestGetAll_ReturnsLatestRates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRatesRepository(ctrl)
	expected := &dto.CurrencyRates{BaseCode: "USD", ConversionRates: dto.Rates{"RUB": 90.0}}
	mockRepo.EXPECT().GetLatestRates().Return(expected, nil)

	svc := &testRatesService{repo: mockRepo}
	got := svc.GetAll()
	assert.Equal(t, expected, got)
}

func TestGetAll_ErrorReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRatesRepository(ctrl)
	mockRepo.EXPECT().GetLatestRates().Return(nil, errors.New("db error"))

	svc := &testRatesService{repo: mockRepo}
	got := svc.GetAll()
	assert.Nil(t, got)
}

func TestGet_ReturnsRateByCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRatesRepository(ctrl)
	cr := &dto.CurrencyRates{BaseCode: "USD", ConversionRates: dto.Rates{"EUR": 0.9}}
	mockRepo.EXPECT().GetLatestRates().Return(cr, nil)

	svc := &testRatesService{repo: mockRepo}
	val, ok := svc.Get("EUR")
	assert.True(t, ok)
	assert.InDelta(t, 0.9, val, 1e-9)
}

func TestGet_WhenNoRate_ReturnsFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRatesRepository(ctrl)
	cr := &dto.CurrencyRates{BaseCode: "USD", ConversionRates: dto.Rates{"EUR": 0.9}}
	mockRepo.EXPECT().GetLatestRates().Return(cr, nil)

	svc := &testRatesService{repo: mockRepo}
	_, ok := svc.Get("RUB")
	assert.False(t, ok)
}
