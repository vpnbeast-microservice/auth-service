package web

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
	logger = commons.GetLogger()
	opts = options.GetAuthServiceOptions()
}

func registerHandlers(router *gin.Engine) {
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":      http.StatusInternalServerError,
				"error":     err,
				"timestamp": time.Now(),
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))
	healthRoutes := router.Group("/health")
	{
		healthRoutes.GET("/ping", pingHandler())
	}
	authRoutes := router.Group("/auth")
	{
		// TODO: single request validator middleware instead of 2 seperate
		authRoutes.POST("/authenticate", authRequestValidator(), authenticateHandler())
		authRoutes.POST("/validate", validateRequestValidator(), validateHandler())
		// TODO: should below /refresh and /whoami endpoints should be GET or POST?
		authRoutes.GET("/refresh", refreshHandler())
		authRoutes.GET("/whoami", whoamiHandler())
	}
}

// InitServer initializes *http.Server with provided parameters
func InitServer(router *gin.Engine) *http.Server {
	registerHandlers(router)
	return &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", opts.ServerPort),
		WriteTimeout: time.Duration(int32(opts.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(opts.ReadTimeoutSeconds)) * time.Second,
	}
}
