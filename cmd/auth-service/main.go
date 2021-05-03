package main

import (
	"auth-service/pkg/config"
	"auth-service/pkg/database"
	"auth-service/pkg/logging"
	"auth-service/pkg/metrics"
	"auth-service/pkg/web"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var (
	serverPort, metricsPort, writeTimeoutSeconds, readTimeoutSeconds int
	logger *zap.Logger
	db *sql.DB
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	logger = logging.GetLogger()
	db = database.GetDatabase()

	// web server/metric server related variables
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

	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	router := gin.Default()
	go metrics.RunMetricsServer(router, metricsPort, writeTimeoutSeconds, readTimeoutSeconds)

	server := web.InitServer(router, fmt.Sprintf(":%d", serverPort), time.Duration(int32(writeTimeoutSeconds)) * time.Second,
		time.Duration(int32(readTimeoutSeconds)) * time.Second)

	logger.Info("web server is up and running", zap.Int("serverPort", serverPort))
	panic(server.ListenAndServe())
}