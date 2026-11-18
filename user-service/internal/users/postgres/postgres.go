package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"lab13/user-service/internal/logger"
	"lab13/user-service/internal/users/service"
)

type querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}
type postgres struct {
	db     querier
	logger *logger.Logger
}

func NewPostgres(db querier, l *logger.Logger) service.IPostgres {
	return &postgres{db: db, logger: l.WithContext("postgres", "users")}
}

func (p postgres) getQuerier(ctx context.Context) querier { return p.db }
