package server

import (
	"log"
	"net/http"

	"go-uniswap-price-fetcher/internal/config"
	"go-uniswap-price-fetcher/internal/service"
)

func Main() {
	config.Init()
	r := service.SetupRouter()
	log.Printf("Listening on %s", config.Cfg.Port)
	if err := http.ListenAndServe(":"+config.Cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
