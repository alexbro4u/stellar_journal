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
	SaveAPOD(apod *nasa_api_models.APODResp) (int64, error)
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
	for {
		apod, err := w.nasaApi.GetAPOD()
		if err != nil {
			w.logger.Error("Failed to get APOD", sl.Err(err))
			time.Sleep(24 * time.Hour)
			continue
		}

		id, err := w.storage.SaveAPOD(apod)
		if err != nil {
			if err == storage.ErrAPODExists {
				w.logger.Info("APOD already exists, retrying in 1 hour")
				time.Sleep(1 * time.Hour)
				continue
			}
			w.logger.Error("Failed to save APOD", sl.Err(err))
		} else {
			w.logger.Info("APOD saved successfully", slog.Int64("id", id))
		}

		time.Sleep(24 * time.Hour)
	}
}
