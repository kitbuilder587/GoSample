package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/kitbuilder587/cryptotrack/internal/db"
)

type mockRepo struct {
	lastLog db.PriceLog
	logs    []db.PriceLog
}

func (m *mockRepo) InsertPrice(ctx context.Context, log db.PriceLog) error {
	m.lastLog = log
	m.logs = append([]db.PriceLog{log}, m.logs...)
	return nil
}
func (m *mockRepo) GetLatest(ctx context.Context, coin string) (db.PriceLog, error) {
	if m.lastLog.Coin == coin {
		return m.lastLog, nil
	}
	return db.PriceLog{}, errors.New("not found")
}
func (m *mockRepo) GetHistory(ctx context.Context, coin string, limit int) ([]db.PriceLog, error) {
	out := []db.PriceLog{}
	for _, l := range m.logs {
		if l.Coin == coin {
			out = append(out, l)
		}
	}
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

type mockClient struct {
	price decimal.Decimal
	err   error
}

func (m *mockClient) GetPriceUSD(ctx context.Context, coin string) (decimal.Decimal, error) {
	return m.price, m.err
}

func TestTrackService_TrackAndSave(t *testing.T) {
	repo := &mockRepo{}
	client := &mockClient{price: decimal.NewFromInt(123)}
	service := NewTrackService(repo, client)

	pl, err := service.TrackAndSave(context.Background(), "bitcoin")
	require.NoError(t, err)
	require.Equal(t, decimal.NewFromInt(123), pl.PriceUSD)
	require.Equal(t, "bitcoin", pl.Coin)
	require.True(t, time.Since(pl.Timestamp) < time.Second)
}

func TestTrackService_TrackAndSave_CoinGeckoErr(t *testing.T) {
	repo := &mockRepo{}
	client := &mockClient{err: errors.New("api fail")}
	service := NewTrackService(repo, client)

	_, err := service.TrackAndSave(context.Background(), "bitcoin")
	require.Error(t, err)
}

func TestTrackService_LatestAndHistory(t *testing.T) {
	repo := &mockRepo{}
	client := &mockClient{price: decimal.NewFromInt(50)}
	service := NewTrackService(repo, client)

	// Сохраняем пару цен
	for i := 0; i < 3; i++ {
		client.price = decimal.NewFromInt(100 + int64(i))
		_, _ = service.TrackAndSave(context.Background(), "bitcoin")
	}
	latest, err := service.Latest(context.Background(), "bitcoin")
	require.NoError(t, err)
	require.True(t, latest.PriceUSD.Equal(decimal.NewFromInt(102)))

	hist, err := service.History(context.Background(), "bitcoin", 2)
	require.NoError(t, err)
	require.Len(t, hist, 2)
}
