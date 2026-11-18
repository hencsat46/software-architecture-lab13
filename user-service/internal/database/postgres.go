package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"lab13/user-service/internal/config"
)

func NewPostgres(ctx context.Context, c config.PostgresConfig) (*pgxpool.Pool, error) {
	p, err := pgxpool.New(ctx, c.DSN())
	if err != nil {
		return nil, err
	}
	if err = p.Ping(ctx); err != nil {
		p.Close()
		return nil, err
	}
	return p, nil
}
