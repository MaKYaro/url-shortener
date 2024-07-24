package router

import (
	"log/slog"
	"net/http"

	"github.com/MaKYaro/url-shortener/internal/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=URLSaver
type URLSaver interface {
	SaveURL(url string) (*domain.Alias, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

type Router struct {
	log    *slog.Logger
	router *http.ServeMux
	saver  URLSaver
	getter URLGetter
}

func NewRouter(
	log *slog.Logger,
	router *http.ServeMux,
	saver URLSaver,
	getter URLGetter,
) *Router {
	return &Router{
		log:    log,
		router: router,
		saver:  saver,
		getter: getter,
	}
}

func (h *Router) InitRoutes() *http.ServeMux {
	h.router.HandleFunc("POST /url", SaveURL(h.log, h.saver))
	h.router.HandleFunc("GET /{alias}", Redirect(h.log, h.getter))
	return h.router
}
