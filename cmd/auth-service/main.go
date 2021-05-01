package main

import (
	"auth-service/pkg/config"
	"auth-service/pkg/logging"
	"auth-service/pkg/metrics"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	serverPort, metricsPort, writeTimeoutSeconds, readTimeoutSeconds int
	logger *zap.Logger
)

func init() {
	logger = logging.GetLogger()
	serverPort = config.GetIntEnv("SERVER_PORT", 5000)
	metricsPort = config.GetIntEnv("METRICS_PORT", 5001)
	writeTimeoutSeconds = config.GetIntEnv("WRITE_TIMEOUT_SECONDS", 10)
	readTimeoutSeconds = config.GetIntEnv("READ_TIMEOUT_SECONDS", 10)
}

func main()  {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()
	go metrics.RunMetricsServer(router, metricsPort, writeTimeoutSeconds, readTimeoutSeconds)

	select {

	}
}