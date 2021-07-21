package options

import (
	"auth-service/pkg/logging"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
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
	err := options.initOptions()
	if err != nil {
		logger.Fatal("fatal error occured while initializing options", zap.Error(err))
	}
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
	PublicKey                  string
	AccessTokenValidInMinutes  int
	RefreshTokenValidInMinutes int
	EncryptionServiceUrl       string
	// database related config
	DbUrl                    string
	DbDriver                 string
	HealthPort               int
	HealthEndpoint           string
	DbMaxOpenConn            int
	DbMaxIdleConn            int
	DbConnMaxLifetimeMin     int
	HealthCheckMaxTimeoutMin int
}

// initOptions initializes AuthServiceOptions while reading environment values, sets default values if not specified
func (aso *AuthServiceOptions) initOptions() error {
	configHost := getStringEnv("CONFIG_SERVER_HOST", "localhost")
	configPort := getIntEnv("CONFIG_SERVER_PORT", 8888)
	appName := getStringEnv("APP_NAME", "auth-service")
	activeProfile := getStringEnv("ACTIVE_PROFILE", "local")
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

	if err := unmarshalConfig("auth-service", aso); err != nil {
		return err
	}

	aso.PrivateKey = strings.Replace(aso.PrivateKey, "\\n", "\n", -1)
	aso.PublicKey = strings.Replace(aso.PublicKey, "\\n", "\n", -1)

	return nil
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