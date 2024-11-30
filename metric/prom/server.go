package prom

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func Serve(ctx context.Context, cfg ServerConfig, logger zerolog.Logger) error {
	logger = logger.With().
		Str("component", "prom http server").
		Logger()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:         cfg.Address(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      mux,
	}

	logger.Info().
		Str("address", cfg.Address()).
		Msg("listening")

	go func() {
		<-ctx.Done()

		closeCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := srv.Shutdown(closeCtx); err != nil {
			logger.Error().
				Err(err).
				Msg("shutdown error")
		}
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
