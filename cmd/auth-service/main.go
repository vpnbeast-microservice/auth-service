package main

import (
	"auth-service/pkg/config"
	"auth-service/pkg/database"
	"auth-service/pkg/logging"
	"auth-service/pkg/metrics"
	"auth-service/pkg/web"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"time"
)

var (
	serverPort, metricsPort, writeTimeoutSeconds, readTimeoutSeconds, dbMaxOpenConn, dbMaxIdleConn,
	dbConnMaxLifetimeMin, healthCheckMaxTimeoutMin int
	dbUrl, dbDriver string
	logger *zap.Logger
)

func init() {
	logger = logging.GetLogger()
	// web server/metric server related variables
	serverPort = config.GetIntEnv("SERVER_PORT", 5000)
	metricsPort = config.GetIntEnv("METRICS_PORT", 5001)
	writeTimeoutSeconds = config.GetIntEnv("WRITE_TIMEOUT_SECONDS", 10)
	readTimeoutSeconds = config.GetIntEnv("READ_TIMEOUT_SECONDS", 10)
	// database related variables
	dbUrl = config.GetStringEnv("DB_URL", "spring:123asd456@tcp(127.0.0.1:3306)/vpnbeast")
	dbDriver = config.GetStringEnv("DB_DRIVER", "mysql")
	dbMaxOpenConn = config.GetIntEnv("DB_MAX_OPEN_CONN", 25)
	dbMaxIdleConn = config.GetIntEnv("DB_MAX_IDLE_CONN", 25)
	dbConnMaxLifetimeMin = config.GetIntEnv("DB_CONN_MAX_LIFETIME_MIN", 5)
	healthCheckMaxTimeoutMin = config.GetIntEnv("HEALTHCHECK_MAX_TIMEOUT_MIN", 5)
}

func main()  {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	db := database.InitDatabase(dbDriver, dbUrl, dbMaxOpenConn, dbMaxIdleConn, dbConnMaxLifetimeMin)

	go func() {
		database.RunHealthProbe(db, healthCheckMaxTimeoutMin)
	}()

	router := mux.NewRouter()
	go metrics.RunMetricsServer(router, metricsPort, writeTimeoutSeconds, readTimeoutSeconds)

	server := web.InitServer(router, fmt.Sprintf(":%d", serverPort), time.Duration(int32(writeTimeoutSeconds)) * time.Second,
		time.Duration(int32(readTimeoutSeconds)) * time.Second)

	logger.Info("web server is up and running", zap.Int("serverPort", serverPort))
	panic(server.ListenAndServe())
}