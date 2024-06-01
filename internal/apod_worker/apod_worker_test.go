package apod_worker_test

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"stellar_journal/internal/apod_worker"
	"stellar_journal/internal/lib/logger/handlers/slogdiscard"
	"stellar_journal/internal/models/nasa_api_models"
	"stellar_journal/internal/storage"
	"testing"
	"time"
)

type MockAPODAPI struct {
	mock.Mock
}

func (m *MockAPODAPI) GetAPOD() (*nasa_api_models.APODResp, error) {
	args := m.Called()
	return args.Get(0).(*nasa_api_models.APODResp), args.Error(1)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveAPOD(apod *nasa_api_models.APODResp) error {
	args := m.Called(apod)
	return args.Error(0)
}

func TestAPODWorkerImpl_Run(t *testing.T) {
	mockAPODAPI := new(MockAPODAPI)
	mockStorage := new(MockStorage)
	logger := slogdiscard.NewDiscardLogger()

	worker := apod_worker.NewAPODWorker(mockAPODAPI, mockStorage, logger)

	t.Run("HappyPath", func(t *testing.T) {
		mockAPODAPI.On("GetAPOD").Return(&nasa_api_models.APODResp{}, nil)
		mockStorage.On("SaveAPOD", mock.Anything).Return(nil)

		go worker.Run()
		time.Sleep(1 * time.Second)
		mockAPODAPI.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("GetAPODFails", func(t *testing.T) {
		mockAPODAPI.On("GetAPOD").Return(nil, errors.New("error"))
		mockStorage.On("SaveAPOD", mock.Anything).Return(nil)

		go worker.Run()
		time.Sleep(1 * time.Second)
		mockAPODAPI.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("SaveAPODFails", func(t *testing.T) {
		mockAPODAPI.On("GetAPOD").Return(&nasa_api_models.APODResp{}, nil)
		mockStorage.On("SaveAPOD", mock.Anything).Return(errors.New("error"))

		go worker.Run()
		time.Sleep(1 * time.Second)
		mockAPODAPI.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("APODAlreadyExists", func(t *testing.T) {
		mockAPODAPI.On("GetAPOD").Return(&nasa_api_models.APODResp{}, nil)
		mockStorage.On("SaveAPOD", mock.Anything).Return(storage.ErrAPODExists)

		go worker.Run()
		time.Sleep(1 * time.Second)
		mockAPODAPI.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}
