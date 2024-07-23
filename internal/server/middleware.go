package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"just_for_fun/pkg/logging"
	"net/http"
	"runtime/debug"
	"time"
)

func AddMiddlewares(r *chi.Mux, logger *logging.DynamicLogger) {
	r.Use(LoggerMiddleware(logger))
	r.Use(CustomRecoverMiddleware(logger))
}

func CustomRecoverMiddleware(logger *logging.DynamicLogger) func(next http.Handler) http.Handler {
	log := logger.GetLogger("access")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					log.Error("Recovered from panic", zap.Any("recovered", rvr), zap.ByteString("stack", debug.Stack()))

					if r.Header.Get("Connection") != "Upgrade" {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func LoggerMiddleware(logger *logging.DynamicLogger) func(next http.Handler) http.Handler {
	logger.AddModule("access")
	log := logger.GetLogger("access")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			msgInfo := fmt.Sprintf("Request %s %d %s in %s", r.Method, ww.Status(), r.URL, time.Since(t1))
			log.Info(msgInfo)
		})
	}
}
