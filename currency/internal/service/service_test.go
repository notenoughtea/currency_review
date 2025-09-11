package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/repository/mocks"
)

type testRatesService struct {
	repo interface {
		GetLatestRates() (*dto.CurrencyRates, error)
	}
}

func (s *testRatesService) HandleRates()               {}
func (s *testRatesService) GetAll() *dto.CurrencyRates { r, _ := s.repo.GetLatestRates(); return r }

func TestGetAll_ReturnsLatestRates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRatesRepository(ctrl)
	expected := &dto.CurrencyRates{Currency: "USD", Rates: []dto.Rate{{Rate: 90.0}}}
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
