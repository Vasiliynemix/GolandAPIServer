package server

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"just_for_fun/internal/config"
	"just_for_fun/internal/server/router"
	"just_for_fun/internal/storage"
	"just_for_fun/pkg"
	"just_for_fun/pkg/logging"
	"net/http"
)

var moduleName = pkg.GetPackageName()

type Server struct {
	logger  *logging.DynamicLogger
	log     zap.Logger
	router  *chi.Mux
	storage *storage.Storage
}

func NewServer(logger *logging.DynamicLogger, router *chi.Mux, storage *storage.Storage) *Server {
	logger.AddModule(moduleName)
	log := logger.GetLogger(moduleName)

	return &Server{
		logger:  logger,
		log:     *log,
		router:  router,
		storage: storage,
	}
}

func (s *Server) Run(cfg *config.Config) {
	AddMiddlewares(s.router, s.logger)

	router.AddRouters(s.router, cfg, s.log, s.storage)

	err := http.ListenAndServe(":"+cfg.Server.Port, s.router)
	if err != nil {
		s.log.Error(err.Error())
		return
	}
}
