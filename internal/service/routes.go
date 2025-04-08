package service

import (
	"github.com/go-chi/chi/v5"

	"go-uniswap-price-fetcher/internal/service/handlers"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/price", handlers.GetPrice)
	return r
}
