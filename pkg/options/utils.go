package options

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"strconv"
)

// getStringEnv gets the specific environment variables with default value, returns default value if variable not set
func getStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// getIntEnv gets the specific environment variables with default value, returns default value if variable not set
func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return convertStringToInt(value)
}

// convertStringToInt converts string environment variables to integer values
func convertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Warn("an error occurred while converting from string to int. Setting it as zero",
			zap.String("error", err.Error()))
		i = 0
	}
	return i
}

func unmarshalConfig(key string, value interface{}) error {
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	bindEnvs(sub)

	return sub.Unmarshal(value)
}

func bindEnvs(sub *viper.Viper) {
	_ = sub.BindEnv("serverPort", "SERVER_PORT")
	_ = sub.BindEnv("metricsPort", "METRICS_PORT")
	_ = sub.BindEnv("metricsEndpoint", "METRICS_ENDPOINT")
	_ = sub.BindEnv("writeTimeoutSeconds", "WRITE_TIMEOUT_SECONDS")
	_ = sub.BindEnv("readTimeoutSeconds", "READ_TIMEOUT_SECONDS")
	_ = sub.BindEnv("issuer", "ISSUER")
	_ = sub.BindEnv("privateKey", "PRIVATE_KEY")
	_ = sub.BindEnv("publicKey", "PUBLIC_KEY")
	_ = sub.BindEnv("accessTokenValidInMinutes", "ACCESS_TOKEN_VALID_IN_MINUTES")
	_ = sub.BindEnv("refreshTokenValidInMinutes", "REFRESH_TOKEN_VALID_IN_MINUTES")
	_ = sub.BindEnv("encryptionServiceUrl", "ENCRYPTION_SERVICE_URL")
	_ = sub.BindEnv("dbUrl", "DB_URL")
	_ = sub.BindEnv("dbDriver", "DB_DRIVER")
	_ = sub.BindEnv("dbMaxOpenConn", "DB_MAX_OPEN_CONN")
	_ = sub.BindEnv("dbMaxIdleConn", "DB_MAX_IDLE_CONN")
	_ = sub.BindEnv("dbConnMaxLifetimeMin", "DB_CONN_MAX_LIFETIME_MIN")
	_ = sub.BindEnv("healthCheckMaxTimeoutMin", "HEALTHCHECK_MAX_TIMEOUT_MIN")
	_ = sub.BindEnv("healthPort", "HEALTH_PORT")
	_ = sub.BindEnv("healthEndpoint", "HEALTH_ENDPOINT")
}
