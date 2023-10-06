package main

import (
	"net/http"
	"os"

	"log/slog"

	"github.com/f1xend/URL_shortener/internal/config"
	"github.com/f1xend/URL_shortener/internal/http-server/handlers/url/save"
	mwlogger "github.com/f1xend/URL_shortener/internal/http-server/middleware/logger"
	"github.com/f1xend/URL_shortener/internal/lib/logger/sl"
	"github.com/f1xend/URL_shortener/internal/lib/logger/sl/handlers/slogpretty"
	"github.com/f1xend/URL_shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shorter", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	_ = storage

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Idle_timeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
	// middleware

	// name, err := storage.GetURL("google")
	// if err != nil {
	// 	log.Error("failed to get URL", sl.Err(err))
	// 	os.Exit(1)
	// }
	// log.Info("get url", slog.StringValue(name))

	// err = storage.DeleteURL("google")
	// if err != nil {
	// 	log.Error("failed to get URL", sl.Err(err))
	// 	os.Exit(1)
	// }

	// name, err := storage.GetURL("google")
	// if err != nil {
	// 	log.Error("failed to get URL", sl.Err(err))
	// 	os.Exit(1)
	// }
	// log.Info("get url", slog.StringValue(name))

	// TODO: run server:

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
