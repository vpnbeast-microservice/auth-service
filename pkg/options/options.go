package options

import (
	"fmt"
	"github.com/spf13/viper"
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var (
	logger *zap.Logger
	opts   *AuthServiceOptions
)

func init() {
	logger = commons.GetLogger()
	opts = newAuthServiceOptions()
	err := opts.initOptions()
	if err != nil {
		logger.Fatal("fatal error occured while initializing options", zap.Error(err))
	}
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

// initOptions initializes AuthServiceOptions while reading environment values, sets default values if not specified
func (aso *AuthServiceOptions) initOptions() error {
	activeProfile := commons.GetStringEnv("ACTIVE_PROFILE", "local")
	appName := commons.GetStringEnv("APP_NAME", "auth-service")
	// TODO: below if/else logic can be implemented using library to decrease duplicate code across other projects?
	if activeProfile == "unit-test" {
		logger.Info("active profile is unit-test, reading configuration from static file")
		// TODO: better approach for that?
		viper.AddConfigPath("./../../config")
		viper.SetConfigName("unit_test")
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	} else {
		configHost := commons.GetStringEnv("CONFIG_SERVER_HOST", "localhost")
		configPort := commons.GetIntEnv("CONFIG_SERVER_PORT", 8888)
		logger.Info("loading configuration from remote server", zap.String("host", configHost),
			zap.Int("port", configPort), zap.String("appName", appName),
			zap.String("activeProfile", activeProfile))
		confAddr := fmt.Sprintf("http://%s:%d/%s-%s.yaml", configHost, configPort, appName, activeProfile)
		resp, err := http.Get(confAddr)
		if err != nil {
			return err
		}

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				panic(err)
			}
		}()

		viper.SetConfigName("application")
		viper.SetConfigType("yaml")
		if err = viper.ReadConfig(resp.Body); err != nil {
			return err
		}
	}

	if err := commons.UnmarshalConfig(appName, aso); err != nil {
		return err
	}

	// required logic for auth-service to convert private key and public key to specific format
	aso.PrivateKey = strings.Replace(aso.PrivateKey, "\\n", "\n", -1)
	aso.PublicKey = strings.Replace(aso.PublicKey, "\\n", "\n", -1)

	return nil
}
