//go:build integration
// +build integration

package coingecko

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestClient_GetPriceUSD_Live(t *testing.T) {
	client := NewClient()

	price, err := client.GetPriceUSD(context.Background(), "bitcoin")
	require.NoError(t, err, "Should not error on real API call")
	require.True(t, price.GreaterThan(decimal.NewFromInt(0)), "Price must be positive")
	t.Logf("BTC price (CoinGecko): %s", price.String())
}

func TestClient_GetPriceUSD_Live_UnknownCoin(t *testing.T) {
	client := NewClient()
	_, err := client.GetPriceUSD(context.Background(), "some_fake_coin_xyz")
	t.Logf("%+v", err)
	require.Error(t, err)
}
