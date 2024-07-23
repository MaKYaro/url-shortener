package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MaKYaro/url-shortener/internal/domain"
	"github.com/MaKYaro/url-shortener/internal/http-server/router/mocks"
	"github.com/MaKYaro/url-shortener/internal/lib/logger/slogdiscard"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		requestBody string
		mockAlias   *domain.Alias
		respError   string
		mockError   error
	}{
		{
			name: "Success",
			url:  "https://pkg.go.dev/net/http#Server.Shutdown",
			mockAlias: &domain.Alias{
				Value:  "asdfasdf",
				URL:    "https://pkg.go.dev/net/http#Server.Shutdown",
				Expire: time.Date(2024, 7, 20, 12, 8, 3, 45, time.UTC),
			},
		},
		{
			name:      "Empty url",
			url:       "",
			mockAlias: &domain.Alias{},
			respError: "invalid request",
		},
		{
			name:        "Failed to decode request body",
			url:         "",
			mockAlias:   &domain.Alias{},
			requestBody: `{"url": "https://ru.wikipedia.org/wiki/Makefile"`,
			respError:   "falied to decode request",
		},
		{
			name: "Mock error",
			url:  "https://www.youtube.com/",
			mockAlias: &domain.Alias{
				Value:  "fasdfasd",
				URL:    "https://www.youtube.com/",
				Expire: time.Date(2024, 7, 20, 12, 8, 3, 45, time.UTC),
			},
			mockError: errors.New("db doesn't work"),
			respError: "can't save url",
		},
	}

	for _, tcase := range cases {
		tcase := tcase
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tcase.respError == "" || tcase.mockError != nil {
				urlSaverMock.
					On("SaveURL", tcase.url).
					Return(tcase.mockAlias, tcase.mockError)
			}

			handler := SaveURL(slogdiscard.NewDiscardLogger(), urlSaverMock)

			var input string
			if tcase.requestBody == "" {
				input = fmt.Sprintf(`{"url": "%s"}`, tcase.url)
			} else {
				input = tcase.requestBody
			}

			req, err := http.NewRequest(
				http.MethodPost,
				"/url",
				bytes.NewReader([]byte(input)),
			)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()
			var resp Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tcase.respError, resp.Error)

		})
	}
}
