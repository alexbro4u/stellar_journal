package all

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "stellar_journal/internal/lib/api/response"
	"stellar_journal/internal/lib/logger/sl"
	"stellar_journal/internal/models/stellar_journal_models"
)

type Response struct {
	resp.Response
	Data []stellar_journal_models.APOD `json:"data"`
}

type APODGetter interface {
	GetJournal() (*[]stellar_journal_models.APOD, error)
}

func New(log *slog.Logger, apodGetter APODGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.journal.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		journals, err := apodGetter.GetJournal()
		if err != nil {
			log.Error("failed to get journals", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get journals"))

			return
		}

		responseOK(w, r, *journals)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data []stellar_journal_models.APOD) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Data:     data,
	})
}
