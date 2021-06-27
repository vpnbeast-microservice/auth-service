package database

import (
	"context"
	"fmt"
	"github.com/etherlabsio/healthcheck"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// RunHealthProbe provides an endpoint, probes the database connection continuously
func RunHealthProbe(router *mux.Router) {
	router.Handle(opts.HealthEndpoint, healthcheck.Handler(
		healthcheck.WithTimeout(time.Duration(int32(opts.HealthCheckMaxTimeoutMin))*time.Second),
		healthcheck.WithChecker(
			"database", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return sqlDB.PingContext(ctx)
				},
			),
		),
	))

	logger.Info("probing mysql", zap.Int("port", opts.HealthPort))
	panic(http.ListenAndServe(fmt.Sprintf(":%d", opts.HealthPort), router))
}
