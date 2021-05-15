package metrics

import (
	"auth-service/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	logger *zap.Logger
)

func init() {
	logger = logging.GetLogger()
}

func RunMetricsServer(router *gin.Engine, metricsPort, writeTimeoutSeconds, readTimeoutSeconds int) {
	metricServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", metricsPort),
		WriteTimeout: time.Duration(int32(writeTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(readTimeoutSeconds)) * time.Second,
	}
	router.GET("/metrics", prometheusHandler())

	logger.Info("metric server is up and running", zap.Int("metricsPort", metricsPort))
	panic(metricServer.ListenAndServe())
}
