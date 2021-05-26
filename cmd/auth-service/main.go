package main

import (
	"auth-service/pkg/database"
	"auth-service/pkg/logging"
	"auth-service/pkg/metrics"
	"auth-service/pkg/options"
	"auth-service/pkg/web"
	"database/sql"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	db     *sql.DB
	opts   *options.AuthServiceOptions
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	logger = logging.GetLogger()
	db = database.GetDatabase()
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
		err := db.Close()
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
