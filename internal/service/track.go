package service

import (
	"context"
	"github.com/kitbuilder587/cryptotrack/internal/db"
	"github.com/shopspring/decimal"
	"time"
)

type PriceClient interface {
	GetPriceUSD(ctx context.Context, coin string) (decimal.Decimal, error)
}

type TrackService struct {
	Repo   db.PriceRepository
	Client PriceClient
}

func NewTrackService(repo db.PriceRepository, client PriceClient) *TrackService {
	return &TrackService{
		Repo:   repo,
		Client: client,
	}
}

func (s *TrackService) TrackAndSave(ctx context.Context, coin string) (db.PriceLog, error) {
	price, err := s.Client.GetPriceUSD(ctx, coin)
	if err != nil {
		return db.PriceLog{}, err
	}
	pl := db.PriceLog{
		Coin:      coin,
		PriceUSD:  price,
		Timestamp: time.Now().UTC(),
	}
	if err := s.Repo.InsertPrice(ctx, pl); err != nil {
		return db.PriceLog{}, err
	}
	return pl, nil
}

func (s *TrackService) Latest(ctx context.Context, coin string) (db.PriceLog, error) {
	return s.Repo.GetLatest(ctx, coin)
}
func (s *TrackService) History(ctx context.Context, coin string, limit int) ([]db.PriceLog, error) {
	return s.Repo.GetHistory(ctx, coin, limit)
}
