package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL:    "https://api.coingecko.com/api/v3",
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type PriceResponse map[string]struct {
	USD decimal.Decimal `json:"usd"`
}

func (c *Client) GetPriceUSD(ctx context.Context, coin string) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", c.BaseURL, coin)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return decimal.Zero, errors.Newf("create request: %w", err)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return decimal.Zero, errors.Newf("request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return decimal.Zero, errors.Newf("unexpected status: %d", resp.StatusCode)
	}
	var pr PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return decimal.Zero, errors.Newf("decode: %w", err)
	}
	coinData, ok := pr[coin]
	if !ok {
		return decimal.Zero, errors.Newf("coin %s not found in response", coin)
	}
	return coinData.USD, nil
}
