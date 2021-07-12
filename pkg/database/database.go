package database

import (
	"auth-service/pkg/logging"
	"auth-service/pkg/options"
	"database/sql"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	logger *zap.Logger
	router *mux.Router
	gormdb *gorm.DB
	sqlDB  *sql.DB
	opts   *options.AuthServiceOptions
	err    error
)

func init() {
	logger = logging.GetLogger()
	router = mux.NewRouter()
	opts = options.GetAuthServiceOptions()
}

// InitDatabase initializes the database connection
func InitDatabase() *gorm.DB {
	gormdb, err = gorm.Open(mysql.Open(opts.DbUrl), &gorm.Config{})
	if err != nil {
		logger.Fatal("fatal error occurred while opening database connection", zap.String("error", err.Error()))
		return nil
	}

	sqlDB, err = gormdb.DB()
	if err != nil {
		logger.Fatal("fatal error occurred while getting sql.DB from gorm.DB", zap.String("error", err.Error()))
		return nil
	}

	tuneDbPooling(sqlDB, opts.DbMaxOpenConn, opts.DbMaxIdleConn, opts.DbConnMaxLifetimeMin)
	go func() {
		RunHealthProbe(router)
	}()

	return gormdb
}

// GetDatabase returns the initialized *sql.DB instance
func GetDatabase() *gorm.DB {
	return gormdb
}

// Read on https://www.alexedwards.net/blog/configuring-sqldb for detailed explanation
func tuneDbPooling(db *sql.DB, dbMaxOpenConn int, dbMaxIdleConn int, dbConnMaxLifetimeMin int) {
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 5. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(dbMaxOpenConn)
	// Set the maximum number of concurrently idle connections to 5. Setting this
	// to less than or equal to 0 will mean that no idle connections are retained.
	db.SetMaxIdleConns(dbMaxIdleConn)
	// Set the maximum lifetime of a connection to 1 hour. Setting it to 0
	// means that there is no maximum lifetime and the connection is reused
	// forever (which is the default behavior).
	db.SetConnMaxLifetime(time.Duration(int32(dbConnMaxLifetimeMin)) * time.Minute)
}
