package handlers

import (
	"math"
	"net/http"
	"strings"

	"go-uniswap-price-fetcher/internal/config"
	"go-uniswap-price-fetcher/internal/service/helpers"
	"go-uniswap-price-fetcher/internal/service/requests"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type PoolPriceResponse struct {
	Token0    string  `json:"token0"`
	Token1    string  `json:"token1"`
	Decimals0 int     `json:"decimals0"`
	Decimals1 int     `json:"decimals1"`
	Price     float64 `json:"price"`
	PriceUSD  *string `json:"priceUSD,omitempty"`
}

func GetPrice(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewPriceRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client, err := ethclient.Dial(config.Cfg.RpcURL)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to connect to rpc").Error(), http.StatusInternalServerError)
		return
	}

	token0, err := helpers.GetTokenData(client, req.Pool, 0)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to fetch first token").Error(), http.StatusInternalServerError)
		return
	}

	token1, err := helpers.GetTokenData(client, req.Pool, 1)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to fetch second token").Error(), http.StatusInternalServerError)
		return
	}

	if token0 == nil || token1 == nil {
		http.Error(w, errors.New("looks like pool not exists").Error(), http.StatusNotFound)
		return
	}

	slot, err := helpers.GetSlot0(client, req.Pool)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to fetch slot0").Error(), http.StatusInternalServerError)
		return
	}

	poolResponse, err := helpers.GetPriceUSD(req.Pool)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to fetch usd price").Error(), http.StatusInternalServerError)
		return
	}

	var priceUSD *string
	if poolResponse != nil {
		if strings.HasPrefix(token0.Symbol, "USD") {
			priceUSD = &poolResponse.Data.Attributes.QuoteTokenPriceUSD
		} else {
			priceUSD = &poolResponse.Data.Attributes.BaseTokenPriceUSD
		}
	}

	price := helpers.CalculatePrice(slot, token0, token1)

	priceResponse := PoolPriceResponse{
		Token0:    token0.Address.String(),
		Token1:    token1.Address.String(),
		Decimals0: int(token0.Decimals),
		Decimals1: int(token1.Decimals),
		Price:     math.Round(price*1e6) / 1e6,
		PriceUSD:  priceUSD,
	}

	helpers.RenderJSON(w, r, priceResponse)
}
