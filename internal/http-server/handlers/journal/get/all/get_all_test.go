package all_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"stellar_journal/internal/models/stellar_journal_models"
	"testing"

	"github.com/stretchr/testify/require"

	"stellar_journal/internal/http-server/handlers/journal/get/all"
	"stellar_journal/internal/http-server/handlers/journal/get/all/mocks"
	"stellar_journal/internal/lib/logger/handlers/slogdiscard"
)

func TestGetAllHandler(t *testing.T) {
	cases := []struct {
		name      string
		respError string
		mockError error
	}{
		{
			name: "Success",
		},
		{
			name:      "GetJournal Error",
			respError: "failed to get journals",
			mockError: errors.New("failed to get journals"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			apodGetterMock := mocks.NewJournalGetter(t)

			if tc.mockError != nil {
				apodGetterMock.On("GetJournal").
					Return(nil, tc.mockError).
					Once()
			} else {
				apodGetterMock.On("GetJournal").
					Return(&[]stellar_journal_models.APOD{}, nil).
					Once()
			}

			handler := all.New(slogdiscard.NewDiscardLogger(), apodGetterMock)

			req, err := http.NewRequest(http.MethodGet, "/journal", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp all.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
