package web

import (
	"auth-service/pkg/config"
	"auth-service/pkg/database"
	"auth-service/pkg/logging"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	logger                                                *zap.Logger
	db                                                    *sql.DB
	accessTokenValidInMinutes, refreshTokenValidInMinutes int
	encryptionServiceUrl                                  string
)

func init() {
	logger = logging.GetLogger()
	db = database.GetDatabase()

	accessTokenValidInMinutes = config.GetIntEnv("ACCESS_TOKEN_VALID_IN_MINUTES", 60)
	refreshTokenValidInMinutes = config.GetIntEnv("REFRESH_TOKEN_VALID_IN_MINUTES", 600)
	encryptionServiceUrl = config.GetStringEnv("ENCRYPTION_SERVICE_URL", "http://localhost:8085/encryption-controller/check")
}

func registerHandlers(router *gin.Engine) {
	router.GET("/health/ping", pingHandler())
	router.POST("/users/authenticate", authenticateHandler())
	// TODO: request validation middleware
	// router.Use(loggingMiddleware)
}

// InitServer initializes *http.Server with provided parameters
func InitServer(router *gin.Engine, serverPort, writeTimeoutSeconds, readTimeoutSeconds int) *http.Server {
	registerHandlers(router)
	return &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", serverPort),
		WriteTimeout: time.Duration(int32(writeTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(readTimeoutSeconds)) * time.Second,
	}
}
