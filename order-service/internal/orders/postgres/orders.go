package postgres

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"lab13/order-service/internal/common/exceptions"
	domain "lab13/order-service/internal/orders"
)

var columns = []string{"id", "user_id", "product", "quantity", "created_at"}

func (p postgres) row(ctx context.Context, query string, args []any) (*domain.Entity, error) {
	l := p.logger.WithMethod("row")
	l.Info("Entering...")
	l.Debug("sql query", zap.String("query", query), zap.Any("args", args))
	var e domain.Entity
	err := p.getQuerier(ctx).QueryRow(ctx, query, args...).Scan(&e.ID, &e.UserID, &e.Product, &e.Quantity, &e.CreatedAt)
	if err != nil {
		l.Error("query failed", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exceptions.ErrNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (p postgres) Create(ctx context.Context, e *domain.Entity) (*domain.Entity, error) {
	l := p.logger.WithMethod("Create")
	l.Info("Entering...")
	q, a, err := sq.
		Insert("orders").
		Columns(columns...).
		Values(e.ID, e.UserID, e.Product, e.Quantity, e.CreatedAt).
		Suffix("RETURNING " + columnNames()).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("build query", zap.Error(err))
		return nil, err
	}
	return p.row(ctx, q, a)
}

func columnNames() string { return "id, user_id, product, quantity, created_at" }

func (p postgres) Fetch(ctx context.Context, id string) (*domain.Entity, error) {
	l := p.logger.WithMethod("Fetch")
	l.Info("Entering...")
	q, a, err := sq.
		Select(columns...).
		From("orders").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("build query", zap.Error(err))
		return nil, err
	}
	return p.row(ctx, q, a)
}

func (p postgres) Update(ctx context.Context, e *domain.Entity) (*domain.Entity, error) {
	l := p.logger.WithMethod("Update")
	l.Info("Entering...")
	b := sq.
		Update("orders")
	b = b.Set("user_id", e.UserID)
	b = b.Set("product", e.Product)
	b = b.Set("quantity", e.Quantity)
	q, a, err := b.
		Where(sq.Eq{"id": e.ID}).
		Suffix("RETURNING " + columnNames()).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("build query", zap.Error(err))
		return nil, err
	}
	return p.row(ctx, q, a)
}

func (p postgres) Delete(ctx context.Context, id string) error {
	l := p.logger.WithMethod("Delete")
	l.Info("Entering...")
	q, a, err := sq.
		Delete("orders").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		l.Error("build query", zap.Error(err))
		return err
	}
	l.Debug("sql query", zap.String("query", q), zap.Any("args", a))
	tag, err := p.getQuerier(ctx).Exec(ctx, q, a...)
	if err != nil {
		l.Error("delete failed", zap.Error(err))
		return err
	}
	if tag.RowsAffected() == 0 {
		l.Warn("entity not found", zap.Error(exceptions.ErrNotFound))
		return exceptions.ErrNotFound
	}
	return nil
}

func (p postgres) List(ctx context.Context) ([]domain.Entity, error) {
	l := p.logger.WithMethod("List")
	l.Info("Entering...")
	q, a, err := sq.
		Select(columns...).
		From("orders").
		OrderBy("created_at", "id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		l.Error("build query", zap.Error(err))
		return nil, err
	}
	l.Debug("sql query", zap.String("query", q), zap.Any("args", a))
	rows, err := p.getQuerier(ctx).Query(ctx, q, a...)
	if err != nil {
		l.Error("query failed", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.Entity, 0)
	for rows.Next() {
		var e domain.Entity
		if err = rows.Scan(&e.ID, &e.UserID, &e.Product, &e.Quantity, &e.CreatedAt); err != nil {
			l.Error("scan failed", zap.Error(err))
			return nil, err
		}
		out = append(out, e)
	}
	if err = rows.Err(); err != nil {
		l.Error("rows failed", zap.Error(err))
		return nil, err
	}
	return out, nil
}
