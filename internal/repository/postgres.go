package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	DbPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	err = DbPool.Ping(ctx)
	if err != nil {
		DbPool.Close()
		return nil, err
	}
	return DbPool, nil
}
