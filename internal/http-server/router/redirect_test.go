package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaKYaro/url-shortener/internal/http-server/router/mocks"
	"github.com/MaKYaro/url-shortener/internal/lib/logger/slogdiscard"
	"github.com/MaKYaro/url-shortener/internal/services"
	"github.com/stretchr/testify/require"
)

func TestRedirect(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		mockURL   string
		mockError error
		respError string
	}{
		{
			name:      "Success",
			alias:     "afddsa",
			mockURL:   "https://ru.wikipedia.org/wiki/Makefile",
			mockError: nil,
		},
		{
			name:      "Alias not found",
			alias:     "dkgkgk",
			mockError: services.ErrAliasNotFound,
			respError: "this url doesn't exist",
		},
		{
			name:      "Can't find url",
			alias:     "kdkdk",
			mockError: services.ErrFailedToFindAlias,
			respError: "can't find url",
		},
	}

	for _, tcase := range cases {
		tcase := tcase
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)
			urlGetterMock.
				On("GetURL", tcase.alias).
				Return(tcase.mockURL, tcase.mockError)

			req, err := http.NewRequest(
				http.MethodGet,
				"/",
				nil,
			)
			require.NoError(t, err)
			req.SetPathValue("alias", tcase.alias)

			redirecter := Redirect(slogdiscard.NewDiscardLogger(), urlGetterMock)
			recorder := httptest.NewRecorder()
			redirecter.ServeHTTP(recorder, req)

			require.Equal(t, tcase.mockURL, recorder.Result().Header.Get("Location"))
		})
	}
}
