package main

import (
	"auth-service/pkg/config"
	"auth-service/pkg/logging"
	"auth-service/pkg/metrics"
	"auth-service/pkg/web"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"time"
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

	server := web.InitServer(router, fmt.Sprintf(":%d", serverPort), time.Duration(int32(writeTimeoutSeconds)) * time.Second,
		time.Duration(int32(readTimeoutSeconds)) * time.Second)

	logger.Info("web server is up and running", zap.Int("serverPort", serverPort))
	panic(server.ListenAndServe())
}