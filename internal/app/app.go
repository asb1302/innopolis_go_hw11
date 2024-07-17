package app

import (
	"github.com/asb1302/innopolis_go_hw11/internal/config"
	"github.com/asb1302/innopolis_go_hw11/internal/handler"
	"github.com/asb1302/innopolis_go_hw11/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.HelloHandler)

	rateLimiter := middleware.RateLimiterMiddleware(cfg.RateLimiter)
	r.Use(rateLimiter)

	return http.ListenAndServe(":8080", r)
}
