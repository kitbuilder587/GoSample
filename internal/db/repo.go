package db

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type PriceRepository interface {
	InsertPrice(ctx context.Context, log PriceLog) error
	GetLatest(ctx context.Context, coin string) (PriceLog, error)
	GetHistory(ctx context.Context, coin string, limit int) ([]PriceLog, error)
}

// Repo wraps an sqlx.DB for convenience
type Repo struct {
	db *sqlx.DB
}

// NewRepo constructs a new Repo using an existing sqlx.DB
func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

// InsertPrice inserts a price log using NamedExec for cleaner bindings
func (r *Repo) InsertPrice(ctx context.Context, log PriceLog) error {
	//goland:noinspection ALL
	_, err := r.db.NamedExecContext(ctx, `
        INSERT INTO price_logs (coin, price_usd, timestamp)
        VALUES (:coin, :price_usd, :timestamp)
    `, log)
	return err
}

// GetLatest retrieves the most recent price for a given coin
func (r *Repo) GetLatest(ctx context.Context, coin string) (PriceLog, error) {
	var log PriceLog
	err := r.db.GetContext(ctx, &log, `
        SELECT id, coin, price_usd, timestamp
        FROM price_logs
        WHERE coin = $1
        ORDER BY timestamp DESC
        LIMIT 1
    `, coin)
	return log, err
}

// GetHistory retrieves the last N price logs for a coin
func (r *Repo) GetHistory(ctx context.Context, coin string, limit int) ([]PriceLog, error) {
	var logs []PriceLog
	err := r.db.SelectContext(ctx, &logs, `
        SELECT id, coin, price_usd, timestamp
        FROM price_logs
        WHERE coin = $1
        ORDER BY timestamp DESC
        LIMIT $2
    `, coin, limit)
	return logs, err
}
