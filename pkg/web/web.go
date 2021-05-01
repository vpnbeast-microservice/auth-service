package web

import (
	"auth-service/pkg/database"
	"auth-service/pkg/logging"
	"database/sql"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	logger *zap.Logger
	db *sql.DB
)


func init() {
	logger = logging.GetLogger()
	db = database.GetDatabase()
}

func registerHandlers(router *mux.Router) {
	pingHandler := http.HandlerFunc(pingHandler)
	router.HandleFunc("/health/ping", pingHandler).Methods("GET").
		Schemes("http").Name("ping")
	router.HandleFunc("/users/authenticate", authenticateHandler).Methods("POST").
		Schemes("http").Name("authenticate")
	// TODO: request validation middleware
	// router.Use(loggingMiddleware)
}

func InitServer(router *mux.Router, addr string, writeTimeout time.Duration, readTimeout time.Duration, ) *http.Server {
	registerHandlers(router)
	return &http.Server{
		Handler: router,
		Addr: addr,
		WriteTimeout: writeTimeout,
		ReadTimeout: readTimeout,
	}
}