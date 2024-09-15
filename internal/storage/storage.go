package storage

import (
	"context"

	"git.codenrock.com/zadanie-6105/internal/storage/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Storage struct {
	l       zerolog.Logger
	pg      *pgxpool.Pool
	Queries *queries.Queries
}

func NewStorage(ctx context.Context, pgConnString string, log zerolog.Logger) (*Storage, error) {
	pgConn, err := pgxpool.New(ctx, pgConnString)
	if err != nil {
		log.Error().Err(err).Msg("failed to init postgres connection")
		return nil, err
	}

	if err = pgConn.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("failed to connect to postgres db")
		pgConn.Close()

		return nil, err
	}

	if err := migration(pgConnString); err != nil {
		log.Error().Err(err).Msg("failed to init migrations")
		return nil, err
	}

	hdl := &Storage{
		l:       log,
		pg:      pgConn,
		Queries: queries.New(pgConn),
	}

	return hdl, nil
}

func (str *Storage) StopPG() {
	if str.pg != nil {
		str.l.Info().Msg("closing PostgreSQL connection pool")
		str.pg.Close()
	}
}

func (str *Storage) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := str.pg.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
