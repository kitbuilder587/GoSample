// Package handlers implements HTTP API for CryptoTrack.
//
//	@title			CryptoTrack API
//	@version		1.0
//	@description	Tracks cryptocurrency prices and stores them in Postgres.
//	@BasePath		/
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kitbuilder587/cryptotrack/internal/service"
)

type API struct {
	Service *service.TrackService
}

func NewAPI(svc *service.TrackService) *API {
	return &API{Service: svc}
}

// Health godoc
// @Summary      Service health check
// @Description  Returns OK if the service is running
// @Tags         health
// @Success      200  {string}  string "ok"
// @Router       /health [get]
func (api *API) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("CryptoTrack is running"))
}

// Track godoc
// @Summary      Track and save coin price
// @Description  Fetches the current price from CoinGecko and saves to DB
// @Tags         tracking
// @Param        coin  query     string  true  "Coin ID"
// @Success      201   {object}  db.PriceLog
// @Failure      400   {string}  string "invalid request"
// @Failure      502   {string}  string "external API error"
// @Router       /track [post]
func (api *API) Track(w http.ResponseWriter, r *http.Request) {
	rawCoin := r.URL.Query().Get("coin")
	coin, err := service.ValidateCoin(rawCoin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pl, err := api.Service.TrackAndSave(r.Context(), coin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(pl)
}

// Latest godoc
// @Summary      Get latest tracked coin price
// @Description  Returns the latest tracked price from DB
// @Tags         query
// @Param        coin  query     string  true  "Coin ID"
// @Success      200   {object}  db.PriceLog
// @Failure      400   {string}  string "invalid request"
// @Failure      404   {string}  string "not found"
// @Router       /latest [get]
func (api *API) Latest(w http.ResponseWriter, r *http.Request) {
	rawCoin := r.URL.Query().Get("coin")
	coin, err := service.ValidateCoin(rawCoin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pl, err := api.Service.Latest(r.Context(), coin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(pl)
}

// History godoc
// @Summary      Get coin price history
// @Description  Returns N latest prices from DB
// @Tags         query
// @Param        coin   query     string  true   "Coin ID"
// @Param        limit  query     int     false  "Limit"  default(10)
// @Success      200    {array}   db.PriceLog
// @Failure      400    {string}  string "invalid request"
// @Failure      500    {string}  string "internal error"
// @Router       /history [get]
func (api *API) History(w http.ResponseWriter, r *http.Request) {
	rawCoin := r.URL.Query().Get("coin")
	coin, err := service.ValidateCoin(rawCoin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := service.ValidateLimit(r.URL.Query().Get("limit"), 10, 1, 1000)
	records, err := api.Service.History(r.Context(), coin, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(records)
}
