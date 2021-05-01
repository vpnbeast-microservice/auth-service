package web

import (
	"auth-service/pkg/logging"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var logger *zap.Logger

func init() {
	logger = logging.GetLogger()
}

func registerHandlers(router *mux.Router) {
	pingHandler := http.HandlerFunc(pingHandler)
	router.HandleFunc("/health/ping", pingHandler).Methods("GET").Schemes("http").Name("ping")
	// router.Use(loggingMiddleware)
}

func InitServer(router *mux.Router, addr string, writeTimeout time.Duration, readTimeout time.Duration) *http.Server {
	registerHandlers(router)
	return &http.Server{
		Handler: router,
		Addr: addr,
		WriteTimeout: writeTimeout,
		ReadTimeout: readTimeout,
	}
}