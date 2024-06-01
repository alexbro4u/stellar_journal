package by_date_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"stellar_journal/internal/models/stellar_journal_models"
	"stellar_journal/internal/storage"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"stellar_journal/internal/http-server/handlers/journal/get/by_date"
	"stellar_journal/internal/http-server/handlers/journal/get/by_date/mocks"
	"stellar_journal/internal/lib/logger/handlers/slogdiscard"
)

func TestGetByDateHandler(t *testing.T) {
	cases := []struct {
		name      string
		date      string
		respError string
		status    int
		mockError error
	}{
		{
			name:   "Success",
			date:   "2022-01-01",
			status: http.StatusOK,
		},
		{
			name:      "GetAPOD Error",
			date:      "2022-01-01",
			respError: "failed to get apod",
			status:    http.StatusInternalServerError,
			mockError: errors.New("failed to get apod"),
		},
		{
			name:      "APOD Not Found",
			date:      "2022-01-01",
			respError: "apod not found",
			status:    http.StatusNotFound,
			mockError: storage.ErrAPODNotFound,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			apodGetterMock := mocks.NewAPODByDateGetter(t)

			if tc.mockError != nil {
				apodGetterMock.On("GetAPOD", mock.Anything).
					Return(nil, tc.mockError).
					Once()
			} else {
				apodGetterMock.On("GetAPOD", mock.Anything).
					Return(&stellar_journal_models.APOD{}, nil).
					Once()
			}

			handler := by_date.New(slogdiscard.NewDiscardLogger(), apodGetterMock)

			router := chi.NewRouter()
			router.Get("/journal/{date}", handler)

			url := fmt.Sprintf("/journal/%s", tc.date)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tc.status, rr.Code)

			body := rr.Body.String()

			var resp by_date.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
