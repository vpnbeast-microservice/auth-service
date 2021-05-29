package options

import (
	"auth-service/pkg/logging"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

var (
	logger  *zap.Logger
	options *AuthServiceOptions
)

func init() {
	logger = logging.GetLogger()
	options = newAuthServiceOptions()
	options.initOptions()
}

// GetAuthServiceOptions returns the initialized AuthServiceOptions
func GetAuthServiceOptions() *AuthServiceOptions {
	return options
}

// newAuthServiceOptions creates an AuthServiceOptions struct with zero values
func newAuthServiceOptions() *AuthServiceOptions {
	return &AuthServiceOptions{}
}

// AuthServiceOptions represents auth-service environment variables
type AuthServiceOptions struct {
	// web server related config
	ServerPort          int
	MetricsPort         int
	MetricsEndpoint     string
	WriteTimeoutSeconds int
	ReadTimeoutSeconds  int
	// jwt related config
	Issuer                     string
	PrivateKey                 string
	AccessTokenValidInMinutes  int
	RefreshTokenValidInMinutes int
	EncryptionServiceUrl       string
	// database related config
	DbUrl                    string
	DbDriver                 string
	HealthPort               int
	DbMaxOpenConn            int
	DbMaxIdleConn            int
	DbConnMaxLifetimeMin     int
	HealthCheckMaxTimeoutMin int
}

// initOptions initializes AuthServiceOptions while reading environment values, sets default values if not specified
func (asc *AuthServiceOptions) initOptions() {
	asc.ServerPort = getIntEnv("SERVER_PORT", 5000)
	asc.MetricsPort = getIntEnv("METRICS_PORT", 5001)
	asc.MetricsEndpoint = getStringEnv("METRICS_ENDPOINT", "/metrics")
	asc.WriteTimeoutSeconds = getIntEnv("WRITE_TIMEOUT_SECONDS", 10)
	asc.ReadTimeoutSeconds = getIntEnv("READ_TIMEOUT_SECONDS", 10)
	asc.Issuer = getStringEnv("ISSUER", "info@thevpnbeast.com")
	asc.PrivateKey = strings.Replace(getStringEnv("PRIVATE_KEY", "-----BEGIN PRIVATE KEY-----\\nMIICdwIBADANBgkqhkiG"+
		"9w0BAQEFAASCAmEwggJdAgEAAoGBANnmLifeLBsiXe/J\\n8O3ophHHaCfJ+EdAUYn7vArJTUtankCD3I8O3n+QM0KNsXzXd+eN6VmNm3bjLp"+
		"Hq\\nVjI/jCr2m1EqXgvRQP74/wOU1sHN3zSRQbcPR0dfJiDfTRmfh/LVrKgcU0kQ4yrG\\nlc0KGB2uslzrKLJCmQ4G0WeM3tKNAgMBAAECgY"+
		"EApsep+FXzSGmLoOfegxqZUe5g\\n6GOMp2yxfH2ztkXR5aVcj2DeRplI8DZ9Jamyei2p1xAl1aevoNXOZV0J0LgXHbm+\\nP6MGU7d+IYD2hI"+
		"CWPfD4pqJafkYc7Q94eQaIiShlYEOoEiLDt09m2V3J/VWxEWw0\\nGTzT1T6zDuwD5epXY/kCQQDv+Xeq+SU5+avfysvm/8bITu/WBRXKxQ7V"+
		"2dg9rJIF\\nrAZSTPUIqdKm2F+o8DIX4sSMouFMgo81Ad4S8D13iQCPAkEA6HNRwmcfQCHzsuuT\\n1407mEPFAcgIckU6e9ubXRRepWPjE6MJ"+
		"IyrDeIkCJfgFPiK8OcNvFLUCD8NaySD1\\nQuHRIwJAJtIqo8QOW6SiQ1/hQItcMwdiETNdZSIf1kSZkNCcBsLfeuzsLuyaIVeb\nkg7Za7fJ"+
		"qB6pZ+EvHZohvNqUdwP4zQJBALf6piiG9C4PcVIYsOA3cYa3hNM/HqhK\\n8NodW9+VAsBGyfC95rqF20aosiGZJ5UhavcRHvc1uNb/GPj99A"+
		"EmuB8CQCy9M89/\\nWGs0V60TrOWn2cmNlvexvxJgtWIjzdtp5rBj/E7Dmfx9nE6sG+uJqob389HYb0fF\\nj8MrN6RCirNhupc=\\n-----END"+
		" PRIVATE KEY-----"), "\\n", "\n", -1)
	asc.AccessTokenValidInMinutes = getIntEnv("ACCESS_TOKEN_VALID_IN_MINUTES", 60)
	asc.RefreshTokenValidInMinutes = getIntEnv("REFRESH_TOKEN_VALID_IN_MINUTES", 600)
	asc.EncryptionServiceUrl = getStringEnv("ENCRYPTION_SERVICE_URL", "http://localhost:8085/encryption-controller/check")
	asc.DbUrl = getStringEnv("DB_URL", "spring:123asd456@tcp(127.0.0.1:3306)/vpnbeast?parseTime=true&loc=Europe%2FIstanbul")
	asc.DbDriver = getStringEnv("DB_DRIVER", "mysql")
	asc.HealthPort = getIntEnv("HEALTH_PORT", 5002)
	asc.DbMaxOpenConn = getIntEnv("DB_MAX_OPEN_CONN", 25)
	asc.DbMaxIdleConn = getIntEnv("DB_MAX_IDLE_CONN", 25)
	asc.DbConnMaxLifetimeMin = getIntEnv("DB_CONN_MAX_LIFETIME_MIN", 5)
	asc.HealthCheckMaxTimeoutMin = getIntEnv("HEALTHCHECK_MAX_TIMEOUT_MIN", 5)
}

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
