package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"syscall"

	"git.codenrock.com/zadanie-6105/config"
	"git.codenrock.com/zadanie-6105/internal/api"
	"git.codenrock.com/zadanie-6105/internal/storage"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:generate sqlc generate

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to init config")
	}

	logLevel, err := zerolog.ParseLevel("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	logger := zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger()

	ctx := context.Background()

	dbLog := logger.With().Str("module", "storage").Logger()

	storage, err := storage.NewStorage(ctx, cfg.PostgresConn, dbLog)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init database")
	}

	apiLog := logger.With().Str("module", "api").Logger()

	APIConfig := &api.Opts{
		Addr:    cfg.ServerAddress,
		Log:     apiLog,
		Storage: storage,
	}

	server, err := api.NewAPI(APIConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start api server")
		}
	}()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info().Msg("awaiting signal")

	sig := <-sigs

	log.Info().Str("signal", sig.String()).Msg("signal received")

	server.Stop()
	storage.StopPG()

	logger.Info().Msg("exiting")
}
