package web

import (
	"auth-service/pkg/database"
	"auth-service/pkg/logging"
	"database/sql"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	logger *zap.Logger
	db *sql.DB
)


func init() {
	logger = logging.GetLogger()
	db = database.GetDatabase()
}

func registerHandlers(router *gin.Engine) {
	router.GET("/health/ping", pingHandler())
	router.POST("/users/authenticate", authenticateHandler())
	// TODO: request validation middleware
	// router.Use(loggingMiddleware)
}

func InitServer(router *gin.Engine, addr string, writeTimeout time.Duration, readTimeout time.Duration, ) *http.Server {
	registerHandlers(router)
	return &http.Server{
		Handler: router,
		Addr: addr,
		WriteTimeout: writeTimeout,
		ReadTimeout: readTimeout,
	}
}