package delete

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/f1xend/URL_shortener/internal/lib/api/response"
	"github.com/f1xend/URL_shortener/internal/lib/logger/sl"
	"github.com/f1xend/URL_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// go:generate go run github.com/vektra/mockery/v2@v2.35.2 --name=URLDeletter
type URLDeletter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeletter URLDeletter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/url/delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is Empty")

			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		err := urlDeletter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("can't delete url", "alias", alias)

			render.JSON(w, r, resp.Error("can't delete"))

			return
		}

		if err != nil {
			log.Info("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("delete url ", slog.String("alias", alias))
	}
}
