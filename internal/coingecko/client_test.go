package coingecko

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestClient_GetPriceUSD_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := &Client{
		BaseURL:    "https://api.coingecko.com/api/v3",
		HTTPClient: &http.Client{},
	}

	httpmock.RegisterResponder("GET",
		"https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd",
		httpmock.NewStringResponder(200, `{"bitcoin":{"usd":67895.12}}`),
	)

	ctx := context.Background()
	price, err := client.GetPriceUSD(ctx, "bitcoin")
	require.NoError(t, err)
	require.True(t, price.Equal(decimal.RequireFromString("67895.12")))
}

func TestClient_GetPriceUSD_CoinNotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := &Client{
		BaseURL:    "https://api.coingecko.com/api/v3",
		HTTPClient: &http.Client{},
	}
	httpmock.RegisterResponder("GET",
		"https://api.coingecko.com/api/v3/simple/price?ids=unknowncoin&vs_currencies=usd",
		httpmock.NewStringResponder(200, `{}`),
	)

	ctx := context.Background()
	_, err := client.GetPriceUSD(ctx, "unknowncoin")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestClient_GetPriceUSD_StatusNot200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := &Client{
		BaseURL:    "https://api.coingecko.com/api/v3",
		HTTPClient: &http.Client{},
	}
	httpmock.RegisterResponder("GET",
		"https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd",
		httpmock.NewStringResponder(429, `rate limit`),
	)

	ctx := context.Background()
	_, err := client.GetPriceUSD(ctx, "bitcoin")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected status")
}
