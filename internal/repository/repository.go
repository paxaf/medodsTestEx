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
	querySetHash string = `INSERT INTO users (guid, refreshToken, useragent) VALUES ($1, $2, $3)`
	queryGetHash string = `SELECT refreshtoken, useragent FROM users WHERE guid = $1`
	queryUpdHash string = `UPDATE users SET refreshToken = $2 WHERE guid = $1`
	queryDelHash string = `DELETE FROM users WHERE guid = $1`
)

type PgRepository interface {
	Close()
	SetHash(guid, hashToken, userAgent string) error
	GetHash(guid string) (string, string, error)
	UpdHash(guid, hashToken string) error
	DeleteHash(guid string) error
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

func (r *Repository) DeleteHash(guid string) error {
	_, err := r.pool.Exec(context.Background(), queryDelHash, guid)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetHash(guid, hashToken, agentHash string) error {
	_, err := r.pool.Exec(context.Background(), querySetHash, guid, hashToken, agentHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdHash(guid, hashToken string) error {
	_, err := r.pool.Exec(context.Background(), queryUpdHash, guid, hashToken)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetHash(guid string) (string, string, error) {
	var hashedtoken, hashedAgent string
	err := r.pool.QueryRow(context.Background(), queryGetHash, guid).Scan(&hashedtoken, &hashedAgent)
	if err != nil {
		return "", "", err
	}
	return hashedtoken, hashedAgent, nil
}
