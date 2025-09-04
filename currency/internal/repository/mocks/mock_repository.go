package mocks

import (
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/notenoughtea/currency_review/currency/internal/dto"
	"github.com/notenoughtea/currency_review/currency/internal/repository"
)

var _ repository.RatesRepository = (*MockRatesRepository)(nil)

type MockRatesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRatesRepositoryMockRecorder
}

type MockRatesRepositoryMockRecorder struct {
	mock *MockRatesRepository
}

func NewMockRatesRepository(ctrl *gomock.Controller) *MockRatesRepository {
	mock := &MockRatesRepository{ctrl: ctrl}
	mock.recorder = &MockRatesRepositoryMockRecorder{mock}
	return mock
}

func (m *MockRatesRepository) EXPECT() *MockRatesRepositoryMockRecorder { return m.recorder }

func (m *MockRatesRepository) StoreRates(newRates dto.CurrencyRates) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreRates", newRates)
	var r0 error
	if rf, ok := ret[0].(func(dto.CurrencyRates) error); ok {
		r0 = rf(newRates)
	} else {
		if ret[0] != nil {
			r0 = ret[0].(error)
		}
	}
	return r0
}

func (mr *MockRatesRepositoryMockRecorder) StoreRates(newRates interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreRates", reflect.TypeOf((*MockRatesRepository)(nil).StoreRates), newRates)
}

func (m *MockRatesRepository) GetLatestRates() (*dto.CurrencyRates, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestRates")
	var r0 *dto.CurrencyRates
	var r1 error
	if rf, ok := ret[0].(func() *dto.CurrencyRates); ok {
		r0 = rf()
	} else {
		if ret[0] != nil {
			r0 = ret[0].(*dto.CurrencyRates)
		}
	}
	if rf, ok := ret[1].(func() error); ok {
		r1 = rf()
	} else {
		if ret[1] != nil {
			r1 = ret[1].(error)
		}
	}
	return r0, r1
}

func (mr *MockRatesRepositoryMockRecorder) GetLatestRates() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestRates", reflect.TypeOf((*MockRatesRepository)(nil).GetLatestRates))
}
