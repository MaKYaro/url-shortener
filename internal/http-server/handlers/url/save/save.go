package save

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MaKYaro/url-shortener/internal/domain"
	resp "github.com/MaKYaro/url-shortener/internal/lib/http/response"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	resp.Response
	Alias  string `json:"alias,omitempty"`
	Expire string `json:"expire,omitempty"`
}

type URLSaver interface {
	SaveURL(url string) (*domain.Alias, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(slog.String("op", op))

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error(
				"failed to decode request body",
				slog.String("error", err.Error()),
			)

			respBytes, _ := json.Marshal(resp.Error("falied to decode request"))
			w.Write(respBytes)

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error(
				"invalid request",
				slog.String("error", err.Error()),
			)
			respBytes, _ := json.Marshal(resp.Error("invalid request"))
			w.Write(respBytes)

			return
		}

		alias, err := urlSaver.SaveURL(req.URL)
		if err != nil {
			log.Error(
				"failed to save url",
				slog.String("error", err.Error()),
			)
			respBytes, _ := json.Marshal(resp.Error("can't save url"))
			w.Write(respBytes)

			return
		}

		log.Info(
			"url added",
			slog.Any("alias value", alias),
		)
		respBytes, _ := json.Marshal(Response{
			Response: resp.OK(),
			Alias:    alias.Value,
			Expire:   alias.ExpireString(),
		})
		w.Write(respBytes)
	}
}
