package by_date

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "stellar_journal/internal/lib/api/response"
	"stellar_journal/internal/lib/logger/sl"
	"stellar_journal/internal/models/stellar_journal_models"
	"stellar_journal/internal/storage"
)

type Response struct {
	resp.Response
	Data stellar_journal_models.APOD `json:"data"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=APODByDateGetter
type APODByDateGetter interface {
	GetAPOD(date string) (*stellar_journal_models.APOD, error)
}

func New(log *slog.Logger, apodGetter APODByDateGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.journal.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		date := chi.URLParam(r, "date")

		apod, err := apodGetter.GetAPOD(date)
		if errors.Is(err, storage.ErrAPODNotFound) {
			log.Error("apod not found", sl.Err(err))

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("apod not found"))

			return
		}
		if err != nil {
			log.Error("failed to get apod", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get apod"))

			return
		}

		responseOK(w, r, *apod)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data stellar_journal_models.APOD) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Data:     data,
	})
}
