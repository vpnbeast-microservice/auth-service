package main

import (
	"auth-service/internal/database"
	"auth-service/internal/metrics"
	"auth-service/internal/options"
	"auth-service/internal/web"
	"github.com/gin-gonic/gin"
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	logger *zap.Logger
	opts   *options.AuthServiceOptions
)

func init() {
	db = database.InitDatabase()
	gin.SetMode(gin.ReleaseMode)
	// gin.DisableConsoleColor()
	logger = commons.GetLogger()
	opts = options.GetAuthServiceOptions()
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
