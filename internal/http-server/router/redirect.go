package router

import (
	"encoding/json"
	"log/slog"
	"net/http"

	resp "github.com/MaKYaro/url-shortener/internal/lib/http/response"
	"github.com/MaKYaro/url-shortener/internal/services"
)

func Redirect(log *slog.Logger, getter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.router.Redirect"

		log.With(slog.String("op", op))

		alias := r.PathValue("alias")

		log.Info("got alias", slog.String("alias", alias))

		url, err := getter.GetURL(alias)
		if err == services.ErrAliasNotFound {
			log.Error(
				"alias not found",
				slog.String("error", err.Error()),
			)
			respByte, _ := json.Marshal(resp.Error("this url doesn't exist"))
			w.Write(respByte)
			return
		}
		if err == services.ErrFailedToFindAlias {
			log.Error(
				"can't find alias",
				slog.String("error", err.Error()),
			)
			respByte, _ := json.Marshal(resp.Error("can't find url"))
			w.Write(respByte)
			return
		}

		log.Info("url successfully found", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)

		log.Info("redirected to", slog.String("url", url))
	}
}
