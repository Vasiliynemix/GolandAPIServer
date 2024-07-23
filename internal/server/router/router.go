package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	_ "just_for_fun/docs"
	"just_for_fun/internal/config"
	"just_for_fun/internal/storage"
	"net/http"
)

type MainRouter struct {
	mr          *chi.Mux
	log         zap.Logger
	BasePattern string
	storage     *storage.Storage
}

func newMainRouter(r *chi.Mux, log zap.Logger, storage *storage.Storage) *MainRouter {
	basePattern := "/api"

	return &MainRouter{mr: r, log: log, BasePattern: basePattern, storage: storage}
}

func AddRouters(r *chi.Mux, cfg *config.Config, log zap.Logger, storage *storage.Storage) {
	mr := newMainRouter(r, log, storage)

	r.Get(cfg.Swagger.Endpoint, httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s%s", cfg.Server.BaseURL, cfg.Swagger.URL)),
	))

	NewUserRouter(mr, storage)

	err := chi.Walk(mr.mr, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Info(fmt.Sprintf("[%s]: '%s' has %d middlewares", method, route, len(middlewares)))
		return nil
	})
	if err != nil {
		log.Error(err.Error())
		return
	}
}
