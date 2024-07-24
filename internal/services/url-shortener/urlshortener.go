package urlshortener

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/MaKYaro/url-shortener/internal/domain"
	"github.com/MaKYaro/url-shortener/internal/storage"
)

var (
	ErrEnableToSave      = errors.New("can't save url")
	ErrAliasNotFound     = errors.New("alias not found")
	ErrFailedToFindAlias = errors.New("failed to find alias")
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=URLSaver
type URLSaver interface {
	SaveURL(alias *domain.Alias) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=URLRemover
type URLRemover interface {
	DeleteURL(alias string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=AliasGenerator
type AliasGenerator interface {
	Generate() string
}

type URLSortener struct {
	log     *slog.Logger
	saver   URLSaver
	getter  URLGetter
	remover URLRemover
	gen     AliasGenerator
	dur     time.Duration
}

func New(
	log *slog.Logger,
	saver URLSaver,
	getter URLGetter,
	remover URLRemover,
	gen AliasGenerator,
	dur time.Duration,
) *URLSortener {
	return &URLSortener{
		log:     log,
		saver:   saver,
		getter:  getter,
		remover: remover,
		gen:     gen,
		dur:     dur,
	}
}

func (u *URLSortener) SaveURL(url string) (*domain.Alias, error) {
	const op = "services.urlshortener.SaveURL"

	log := u.log.With(slog.String("op", op))

	expire := time.Now().Add(u.dur)
	aliasValue := u.gen.Generate()
	aliasToSave := domain.Alias{Value: aliasValue, URL: url, Expire: expire}
	err := u.saver.SaveURL(&aliasToSave)

	for err == storage.ErrAliasExists {
		log.Info(
			"alias already exists, try to generate unique",
			slog.String("alias", aliasValue),
		)
		aliasValue = u.gen.Generate()
		aliasToSave = domain.Alias{Value: aliasValue, URL: url, Expire: expire}
		err = u.saver.SaveURL(&aliasToSave)
	}

	if err != nil {
		log.Error(
			"can't save alias",
			slog.String("err", err.Error()),
		)
		return nil, fmt.Errorf("%s: %w", op, ErrEnableToSave)
	}

	return &aliasToSave, nil
}

func (u *URLSortener) GetURL(alias string) (string, error) {
	const op = "services.urlshortener.GetURL"

	log := u.log.With(slog.String("op", op))

	url, err := u.getter.GetURL(alias)
	if err == storage.ErrURLNotFound {
		log.Error(
			"alias not found",
			slog.String("error", err.Error()),
		)
		return "", fmt.Errorf("%s: %w", op, ErrAliasNotFound)
	}

	if err != nil {
		log.Error(
			"failed to find alias",
			slog.String("error", err.Error()),
		)
		return "", fmt.Errorf("%s: %w", op, ErrFailedToFindAlias)
	}

	return url, nil
}
