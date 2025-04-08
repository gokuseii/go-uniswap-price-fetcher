package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type PoolResponse struct {
	Data struct {
		Attributes struct {
			BaseTokenPriceUSD  string `json:"base_token_price_usd"`
			QuoteTokenPriceUSD string `json:"quote_token_price_usd"`
		} `json:"attributes"`
	} `json:"data"`
}

func GetPriceUSD(poolAddress string) (*PoolResponse, error) {
	url := fmt.Sprintf("https://api.geckoterminal.com/api/v2/networks/eth/pools/%s", poolAddress)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	var result PoolResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &result, nil
}
