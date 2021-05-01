package database

import (
	"auth-service/pkg/config"
	"auth-service/pkg/logging"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"time"
)

var (
	logger *zap.Logger
	dbMaxOpenConn, dbMaxIdleConn, dbConnMaxLifetimeMin, healthCheckMaxTimeoutMin, healthPort int
	dbUrl, dbDriver string
	router *mux.Router
	db *sql.DB
)

func init() {
	logger = logging.GetLogger()
	router = mux.NewRouter()

	// database related variables
	dbUrl = config.GetStringEnv("DB_URL", "spring:123asd456@tcp(127.0.0.1:3306)/vpnbeast?parseTime=true")
	dbDriver = config.GetStringEnv("DB_DRIVER", "mysql")
	healthPort = config.GetIntEnv("HEALTH_PORT", 5002)
	dbMaxOpenConn = config.GetIntEnv("DB_MAX_OPEN_CONN", 25)
	dbMaxIdleConn = config.GetIntEnv("DB_MAX_IDLE_CONN", 25)
	dbConnMaxLifetimeMin = config.GetIntEnv("DB_CONN_MAX_LIFETIME_MIN", 5)
	healthCheckMaxTimeoutMin = config.GetIntEnv("HEALTHCHECK_MAX_TIMEOUT_MIN", 5)

	db = initDatabase(dbDriver, dbUrl, dbMaxOpenConn, dbMaxIdleConn, dbConnMaxLifetimeMin)
}

func initDatabase(dbDriver, dbUrl string, dbMaxOpenConn, dbMaxIdleConn, dbConnMaxLifetimeMin int) *sql.DB {
	db, err := sql.Open(dbDriver, dbUrl)
	if err != nil {
		logger.Fatal("fatal error occured while opening database connection", zap.String("error", err.Error()))
	}
	tuneDbPooling(db, dbMaxOpenConn, dbMaxIdleConn, dbConnMaxLifetimeMin)

	go func() {
		RunHealthProbe(router, db, healthCheckMaxTimeoutMin, healthPort)
	}()

	return db
}

func GetDatabase() *sql.DB {
	return db
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