package options

import (
	"auth-service/pkg/logging"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
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

	if err := unmarshalConfig(appName, aso); err != nil {
		return err
	}

	// required logic for auth-service to convert private key and public key to specific format
	aso.PrivateKey = strings.Replace(aso.PrivateKey, "\\n", "\n", -1)
	aso.PublicKey = strings.Replace(aso.PublicKey, "\\n", "\n", -1)

	return nil
}
