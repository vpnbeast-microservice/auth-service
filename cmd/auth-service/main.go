package main

import (
	"auth-service/pkg/database"
	"auth-service/pkg/logging"
	"auth-service/pkg/metrics"
	"auth-service/pkg/options"
	"auth-service/pkg/web"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	logger *zap.Logger
	db     *gorm.DB
	opts   *options.AuthServiceOptions
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	// gin.DisableConsoleColor()
	logger = logging.GetLogger()
	opts = options.GetAuthServiceOptions()
	db = database.GetDatabase()
}

func main() {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	defer func() {
		sqlDb, err := db.DB()
		if err != nil {
			panic(err)
		}

		err = sqlDb.Close()
		if err != nil {
			panic(err)
		}
	}()

	router := gin.Default()
	go metrics.RunMetricsServer(router)
	server := web.InitServer(router)
	logger.Info("web server is up and running", zap.Int("serverPort", opts.ServerPort))
	panic(server.ListenAndServe())
}
