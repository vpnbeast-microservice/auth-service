package database

import (
	"fmt"
	"github.com/dimiro1/health"
	"github.com/dimiro1/health/db"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

// RunHealthProbe provides an endpoint, probes the database connection continuously
func RunHealthProbe(router *mux.Router) {
	mysql := db.NewMySQLChecker(sqlDB)
	handler := health.NewHandler()
	handler.AddChecker("MySQL", mysql)
	router.Handle(opts.HealthEndpoint, handler)
	logger.Info("probing mysql", zap.Int("port", opts.HealthPort))
	panic(http.ListenAndServe(fmt.Sprintf(":%d", opts.HealthPort), router))
}
