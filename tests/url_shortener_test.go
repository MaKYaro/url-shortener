package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/MaKYaro/url-shortener/internal/http-server/router"
	"github.com/MaKYaro/url-shortener/internal/lib/api"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	url := url.URL{
		Scheme: "http",
		Host:   host,
	}
	expected := httpexpect.Default(t, url.String())

	expected.POST("/url").
		WithJSON(router.Request{
			URL: gofakeit.URL(),
		}).
		Expect().Status(200).
		JSON().Object().ContainsKey("alias")

}

func TestURLShortener_SaveRedirect(t *testing.T) {
	cases := []struct {
		name  string
		url   string
		error string
	}{
		{
			name: "Valid URL",
			url:  gofakeit.URL(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid",
			error: "invalid request",
		},
	}
	for _, tcase := range cases {
		tcase := tcase
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			url := url.URL{
				Scheme: "http",
				Host:   host,
			}

			expected := httpexpect.Default(t, url.String())

			// Save
			resp := expected.POST("/url").
				WithJSON(router.Request{
					URL: tcase.url,
				}).
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tcase.error != "" {
				resp.NotContainsKey("alias")
				resp.Value("error").String().IsEqual(tcase.error)
				return
			}

			resp.Value("alias").String().NotEmpty()
			alias := resp.Value("alias").String().Raw()

			// Redirect
			testRedirect(t, alias, tcase.url)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}
