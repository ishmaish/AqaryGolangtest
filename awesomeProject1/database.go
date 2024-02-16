package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func NewDBPool(connString string) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &DB{pool}, nil
}

// Implement database operations here
