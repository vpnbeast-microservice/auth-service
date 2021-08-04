package options

import (
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"strings"
)

var (
	logger *zap.Logger
	opts   *AuthServiceOptions
)

func init() {
	logger = commons.GetLogger()
	opts = newAuthServiceOptions()
	err := commons.InitOptions(opts, "auth-service")
	if err != nil {
		logger.Fatal("fatal error occured while initializing options", zap.Error(err))
	}
	// required logic for auth-service to convert private key and public key to specific format
	opts.PrivateKey = strings.Replace(opts.PrivateKey, "\\n", "\n", -1)
	opts.PublicKey = strings.Replace(opts.PublicKey, "\\n", "\n", -1)
}

// GetAuthServiceOptions returns the initialized AuthServiceOptions
func GetAuthServiceOptions() *AuthServiceOptions {
	return opts
}

// newAuthServiceOptions creates an AuthServiceOptions struct with zero values
func newAuthServiceOptions() *AuthServiceOptions {
	return &AuthServiceOptions{}
}

// AuthServiceOptions represents auth-service environment variables
type AuthServiceOptions struct {
	// web server related config
	ServerPort          int    `env:"SERVER_PORT"`
	MetricsPort         int    `env:"METRICS_PORT"`
	MetricsEndpoint     string `env:"METRICS_ENDPOINT"`
	WriteTimeoutSeconds int    `env:"WRITE_TIMEOUT_SECONDS"`
	ReadTimeoutSeconds  int    `env:"READ_TIMEOUT_SECONDS"`
	// jwt related config
	Issuer                     string `env:"ISSUER"`
	PrivateKey                 string `env:"PRIVATE_KEY"`
	PublicKey                  string `env:"PUBLIC_KEY"`
	AccessTokenValidInMinutes  int    `env:"ACCESS_TOKEN_VALID_IN_MINUTES"`
	RefreshTokenValidInMinutes int    `env:"REFRESH_TOKEN_VALID_IN_MINUTES"`
	EncryptionServiceUrl       string `env:"ENCRYPTION_SERVICE_URL"`
	// database related config
	DbUrl                    string `env:"DB_URL"`
	DbDriver                 string `env:"DB_DRIVER"`
	HealthPort               int    `env:"HEALTH_PORT"`
	HealthEndpoint           string `env:"HEALTH_ENDPOINT"`
	DbMaxOpenConn            int    `env:"DB_MAX_OPEN_CONN"`
	DbMaxIdleConn            int    `env:"DB_MAX_IDLE_CONN"`
	DbConnMaxLifetimeMin     int    `env:"DB_CONN_MAX_LIFETIME_MIN"`
	HealthCheckMaxTimeoutMin int    `env:"HEALTHCHECK_MAX_TIMEOUT_MIN"`
}