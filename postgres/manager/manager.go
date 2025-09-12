package manager

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKey string

const querierKey ctxKey = "querier"

func QuerierFromCtx[T any](ctx context.Context) (T, bool) {
	v := ctx.Value(querierKey)
	if v == nil {
		var zero T
		return zero, false
	}
	return v.(T), true
}

func ContextWithQuerier[T any](ctx context.Context, querier T) context.Context {
	return context.WithValue(ctx, querierKey, querier)
}

type Manager[T any] struct {
	pool    *pgxpool.Pool
	factory func(pgx.Tx) T
}

func New[T any](pool *pgxpool.Pool, factory func(pgx.Tx) T) *Manager[T] {
	return &Manager[T]{pool: pool, factory: factory}
}

func (m *Manager[T]) Do(ctx context.Context, f func(context.Context) error) error {
	// Уже есть querier в контексте
	if _, ok := QuerierFromCtx[T](ctx); ok {
		return f(ctx)
	}

	// Начинаем новую транзакцию
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	querier := m.factory(tx)
	ctx = ContextWithQuerier(ctx, querier)

	if err = f(ctx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	tx = nil
	return nil
}
