package database

import (
	"context"
	"database/sql"
	"github.com/etherlabsio/healthcheck"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RunHealthProbe(db *sql.DB, healthCheckMaxTimeoutMin int) {
	router := mux.NewRouter()
	router.Handle("/health", healthcheck.Handler(
		healthcheck.WithTimeout(time.Duration(int32(healthCheckMaxTimeoutMin)) * time.Second),
		healthcheck.WithChecker(
			"database", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return db.PingContext(ctx)
				},
			),
		),
	))

	logger.Info("metric server is up and running", zap.Int("port", 9290))
	panic(http.ListenAndServe(":9290" , router))
}