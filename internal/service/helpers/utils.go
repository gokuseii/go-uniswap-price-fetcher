package helpers

import (
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/pkg/errors"
)

func RenderJSON(w http.ResponseWriter, r *http.Request, data any) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(errors.Wrap(err, "failed to render response"))
	}
}

func CalculatePrice(slot *Slot0, token0, token1 *Token) float64 {
	sqrtPriceX96 := slot.SqrtPriceX96
	decimals0 := int64(token0.Decimals)
	decimals1 := int64(token1.Decimals)
	fSqrtPriceX96 := new(big.Float).SetInt(sqrtPriceX96)
	twoPow96 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96))

	// (sqrtPriceX96 / 2^96)^2
	priceFloat := new(big.Float).Quo(fSqrtPriceX96, twoPow96)
	priceFloat.Mul(priceFloat, priceFloat)

	// (10^decimals0 / 10^decimals1)
	ten := big.NewInt(10)
	pow0 := new(big.Int).Exp(ten, big.NewInt(int64(decimals0)), nil)
	pow1 := new(big.Int).Exp(ten, big.NewInt(int64(decimals1)), nil)
	scale0 := new(big.Float).SetInt(pow0)
	scale1 := new(big.Float).SetInt(pow1)
	scaleRatio := new(big.Float).Quo(scale0, scale1)
	priceFloat.Mul(priceFloat, scaleRatio)

	result, _ := priceFloat.Float64()
	return result
}
