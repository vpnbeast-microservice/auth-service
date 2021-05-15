package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/etherlabsio/healthcheck"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RunHealthProbe(router *mux.Router, db *sql.DB, healthCheckMaxTimeoutMin, healthPort int) {
	router.Handle("/health", healthcheck.Handler(
		healthcheck.WithTimeout(time.Duration(int32(healthCheckMaxTimeoutMin))*time.Second),
		healthcheck.WithChecker(
			"database", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return db.PingContext(ctx)
				},
			),
		),
	))

	logger.Info("probing mysql", zap.Int("port", healthPort))
	panic(http.ListenAndServe(fmt.Sprintf(":%d", healthPort), router))
}
