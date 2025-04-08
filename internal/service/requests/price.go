package requests

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gorilla/schema"
)

type PriceRequest struct {
	Pool string `schema:"pool" validate:"required,validPool"`
}

func NewPriceRequest(r *http.Request) (PriceRequest, error) {
	var req PriceRequest
	decoder := schema.NewDecoder()
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		return req, err
	}

	if err := req.Validate(); err != nil {
		return req, err
	}
	return req, nil
}

func (r *PriceRequest) Validate() error {
	if r.Pool == "" {
		return errors.New("pool is required")
	}

	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(r.Pool) {
		return errors.New("invalid pool address format")
	}

	return nil
}
