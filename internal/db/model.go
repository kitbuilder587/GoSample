package db

import (
	"github.com/shopspring/decimal"
	"time"
)

// PriceLog represents a record in price_logs table.
//
// swagger:model
type PriceLog struct {
	ID        int64           `db:"id"`
	Coin      string          `db:"coin"`
	PriceUSD  decimal.Decimal `db:"price_usd"`
	Timestamp time.Time       `db:"timestamp"`
}
