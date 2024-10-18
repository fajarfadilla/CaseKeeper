package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/fajarfadilla/casekeeper/backend/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func gracefulShutdown(apiServer *http.Server) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}

func serverStart() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.yml", "Config file name")

	flag.Parse()

	cfg := defaultConfig()
	cfg.loadFromEnv()

	if len(configFile) > 0 {
		err := loadConfigFromFile(configFile, &cfg)
		if err != nil {
			log.Warn().Str("file", configFile).Err(err).Msg("Cannot load config file, use defaults")
		}
	}

	log.Debug().Any("config", cfg).Msg("Config loaded")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBConfig.ConnStr())
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
	}

	user.SetPool(pool)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Mount("/api/v1/user/", user.Router())

	log.Info().Msgf("Starting up a server on %s", cfg.Listen.Addr())

	srv := &http.Server{
		Addr:    ":6969",
		Handler: r,
	}
	go gracefulShutdown(srv)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
	log.Info().Msg("server stopped")
}
