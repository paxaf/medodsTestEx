package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFeedbackExists = errors.New("feedback already exists")
	ErrPostgreSQL     = errors.New("PostgreSQL error")
)

const (
	querySetHash string = `INSERT INTO users (guid, refreshToken) VALUES ($1, $2)`
)

type PgRepository interface {
	Close()
	SetHash(guid, hashToken string) error
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) PgRepository {
	return &Repository{pool: pool}
}

func (r *Repository) Close() {
	r.pool.Close()
}

func (r *Repository) SetHash(guid, hashToken string) error {
	_, err := r.pool.Exec(context.Background(), querySetHash, guid, hashToken)
	if err != nil {
		return err
	}
	return nil
}
