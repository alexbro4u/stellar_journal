package apod_worker

import (
	"log/slog"
	"stellar_journal/internal/lib/logger/sl"
	"stellar_journal/internal/models/nasa_api_models"
	"stellar_journal/internal/storage"
	"time"
)

type APODAPI interface {
	GetAPOD() (*nasa_api_models.APODResp, error)
}

type Storage interface {
	SaveAPOD(apod *nasa_api_models.APODResp) error
}

type APODWorkerImpl struct {
	nasaApi APODAPI
	storage Storage
	logger  *slog.Logger
}

func NewAPODWorker(nasaApi APODAPI, storage Storage, logger *slog.Logger) *APODWorkerImpl {
	return &APODWorkerImpl{
		nasaApi: nasaApi,
		storage: storage,
		logger:  logger,
	}
}

func (w *APODWorkerImpl) Run() {
	errCount := 0
	for {
		apod, err := w.nasaApi.GetAPOD()
		if err != nil {
			w.logger.Error("Failed to get APOD", sl.Err(err))
			time.Sleep(1 * time.Hour)
			continue
		}

		err = w.storage.SaveAPOD(apod)
		if err == storage.ErrAPODExists {
			errCount++
			waitTime := 1 * time.Hour
			if errCount >= 2 {
				waitTime = 24 * time.Hour
			}
			w.logger.Info("APOD already exists, retrying after wait time")
			time.Sleep(waitTime)
			continue
		}

		if err != nil {
			w.logger.Error("Failed to save APOD", sl.Err(err))
		} else {
			w.logger.Info("APOD saved successfully")
		}

		errCount = 0
		time.Sleep(24 * time.Hour)
	}
}
