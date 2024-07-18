package urlshortener

import (
	"errors"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/MaKYaro/url-shortener/internal/domain"
	"github.com/MaKYaro/url-shortener/internal/lib/logger/slogdiscard"
	"github.com/MaKYaro/url-shortener/internal/services/url-shortener/mocks"
	"github.com/MaKYaro/url-shortener/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestSaveURL(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2024, 7, 18, 16, 5, 7, 50, time.UTC)
	})
	now := time.Now()
	dur, _ := time.ParseDuration("3h")
	expire := now.Add(dur)

	cases := []struct {
		name           string
		url            string
		aliasToSave    *domain.Alias
		generatedAlias string
		mockError      error
		saveError      error
	}{
		{
			name: "Success",
			url:  "http://github.com/stretchr/testify/assert",
			aliasToSave: &domain.Alias{
				Value:  "adf3KY7r",
				URL:    "http://github.com/stretchr/testify/assert",
				Expire: expire,
			},
			generatedAlias: "adf3KY7r",
			mockError:      nil,
			saveError:      nil,
		},
		{
			name: "SaveURL error",
			url:  "http://github.com/stretchr/testify/assert",
			aliasToSave: &domain.Alias{
				Value:  "adf3KY7r",
				URL:    "http://github.com/stretchr/testify/assert",
				Expire: expire,
			},
			generatedAlias: "adf3KY7r",
			mockError:      errors.New("db doesn't work"),
			saveError:      ErrEnableToSave,
		},
		{
			name: "Empty url",
			url:  "",
			aliasToSave: &domain.Alias{
				Value:  "adf3KY7r",
				URL:    "",
				Expire: expire,
			},
			generatedAlias: "adf3KY7r",
			mockError:      errors.New("empty url"),
			saveError:      ErrEnableToSave,
		},
	}

	for _, tcase := range cases {
		tcase := tcase
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)
			urlSaverMock.On("SaveURL", tcase.aliasToSave).Return(tcase.mockError).Once()

			aliasGenerator := mocks.NewAliasGenerator(t)
			aliasGenerator.On("Generate").Return(tcase.generatedAlias).Once()

			urlShortener := New(
				slogdiscard.NewDiscardLogger(),
				urlSaverMock,
				nil,
				nil,
				aliasGenerator,
				dur,
			)

			savedAlias, err := urlShortener.SaveURL(tcase.url)

			if tcase.saveError != nil {
				require.Nil(t, savedAlias)
				require.ErrorIs(t, err, tcase.saveError)
			} else {
				require.Equal(t, savedAlias, tcase.aliasToSave)
				require.ErrorIs(t, err, tcase.saveError)
			}
		})
	}
}

func TestSaveURLErrAliasExists(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2024, 7, 18, 16, 5, 7, 50, time.UTC)
	})

	now := time.Now()
	dur, _ := time.ParseDuration("3h")
	expire := now.Add(dur)
	aliasToSave := &domain.Alias{
		Value:  "adf3KY7r",
		URL:    "http://github.com/stretchr/testify/assert",
		Expire: expire,
	}

	urlSaverMock := mocks.NewURLSaver(t)
	urlSaverMock.On("SaveURL", aliasToSave).Return(storage.ErrAliasExists).Times(6)
	urlSaverMock.On("SaveURL", aliasToSave).Return(nil).Once()

	aliasGenerator := mocks.NewAliasGenerator(t)
	aliasGenerator.On("Generate").Return("adf3KY7r")

	urlShortener := New(
		slogdiscard.NewDiscardLogger(),
		urlSaverMock,
		nil,
		nil,
		aliasGenerator,
		dur,
	)

	savedAlias, err := urlShortener.SaveURL("http://github.com/stretchr/testify/assert")

	require.Nil(t, err)
	require.Equal(t, savedAlias, aliasToSave)
}

func TestGetURL(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		mockURL   string
		mockError error
		wantURL   string
		wantError error
	}{
		{
			name:      "Success",
			alias:     "sfaa9rl",
			mockURL:   "http://github.com/stretchr/testify/assert",
			mockError: nil,
			wantURL:   "http://github.com/stretchr/testify/assert",
			wantError: nil,
		},
		{
			name:      "URL not found error ",
			alias:     "sfaa9rl",
			mockURL:   "",
			mockError: storage.ErrURLNotFound,
			wantURL:   "",
			wantError: ErrURLNotFound,
		},
		{
			name:      "Can't find url error",
			alias:     "sfaa9rl",
			mockURL:   "",
			mockError: errors.New("can't find url for some reasons"),
			wantURL:   "",
			wantError: ErrCantFindUrl,
		},
	}

	for _, tcase := range cases {
		tcase := tcase
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)
			urlGetterMock.On("GetURL", tcase.alias).
				Return(tcase.mockURL, tcase.mockError).
				Once()

			dur, _ := time.ParseDuration("3s")
			urlShortener := New(
				slogdiscard.NewDiscardLogger(),
				nil,
				urlGetterMock,
				nil,
				nil,
				dur,
			)

			gotURL, gotError := urlShortener.GetURL(tcase.alias)

			require.Equal(t, tcase.wantURL, gotURL)
			require.ErrorIs(t, gotError, tcase.wantError)

		})
	}
}
