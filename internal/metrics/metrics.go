package metrics

import (
	"auth-service/internal/options"
	"fmt"
	"github.com/gin-gonic/gin"
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	logger *zap.Logger
	opts   *options.AuthServiceOptions
)

func init() {
	opts = options.GetAuthServiceOptions()
	logger = commons.GetLogger()
}

// RunMetricsServer provides an endpoint, exports prometheus metrics using prometheus client golang
func RunMetricsServer(router *gin.Engine) {
	metricServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", opts.MetricsPort),
		WriteTimeout: time.Duration(int32(opts.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(opts.ReadTimeoutSeconds)) * time.Second,
	}
	router.GET(opts.MetricsEndpoint, prometheusHandler())

	logger.Info("metric server is up and running", zap.Int("metricsPort", opts.MetricsPort))
	panic(metricServer.ListenAndServe())
}
