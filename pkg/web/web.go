package web

import (
	"auth-service/pkg/database"
	"auth-service/pkg/logging"
	"auth-service/pkg/options"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	logger *zap.Logger
	db     *sql.DB
	opts   *options.AuthServiceOptions
)

func init() {
	logger = logging.GetLogger()
	db = database.GetDatabase()
	opts = options.GetAuthServiceOptions()
}

func registerHandlers(router *gin.Engine) {
	router.GET("/health/ping", pingHandler())
	router.POST("/users/authenticate", authenticateHandler())
	// TODO: request validation middleware
	// router.Use(loggingMiddleware)
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
