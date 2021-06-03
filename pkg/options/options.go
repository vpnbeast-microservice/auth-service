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
	PublicKey				   string
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
	asc.PrivateKey = strings.Replace(getStringEnv("PRIVATE_KEY", "-----BEGIN PRIVATE KEY-----\\nMIIEvQI"+
		"BADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDAqVXTwDGuM+87\\npf9Oz632tZQ2j6jhkSoNEYHYlhTqBNMEP/N/6dqwtJzwKCBRG"+
		"TC2N8+BUJbynx6x\\nXMaCmjLQ+tRSYMfcwf0UoZWspl8uDCB/OY4NZweLXhFEuHMTV/u8M6bZhIIKu/9R\\naCaLylJgXDXhTZ8iFMP5YTrH"+
		"EFXxVowMVGNCYKrlWteJ7rTydJq/1Jo6muV2nXOX\\nDSm5KrOJu2sVLpjXIeehupK5n7PrcK/Dr1RbisihvKiPu+14YEdrdfP7r4+KuVha\\n"+
		"x9PLB7YRR2S3ttuJckkAwEuevx1hPWn3rtspg8Zi6fu+65Llk7jTHPHexrMaXd8n\\n7qydWMP7AgMBAAECggEBAK4bGTXP1NWhl0tnOq6lHY"+
		"f7FeKstqiJv9+pd5ccIfBN\\nxchsZoes9PWFcuHQ0UuYoR26l+o7wv2k8F7GAZra8JtMYX3Eejk0kZoIYDNy8kax\\nrNhsUcQo3HeE3cQlj"+
		"9DmTNcKMnkVt1MuC5Ast9DSWNk9228s72ckLun5hN2KFLCP\\nudPM5jERYEwO0PoeUM9nY/8kx+JLhCXkvthwIlB1pBPZDfyFcG5qxT+dmTQ"+
		"I4Pg5\\nNl/SFJHX2ALbcN0WStEoe/FuLT/08lZgb5VADpvH+KKvSaMihjlkxzwrJT0cAuss\\nnbx9ub3x0VOjb5pRtbry/PhkLLwhiZD1Q"+
		"37hXGltqoECgYEA/j4KxymxUdufgunC\\nUPs5pDdJi10ZGlU3aDOTgF/7VtJuJk+KOat4/Q8zCzk+rL6+wThkUekBUpJO8YM8\\nWNLMbVs"+
		"Pq7aD8uZjy1T9UQ2xClHvYRWVjPsRHJtgAzTp6FsYBJxBodpZB+6LCL35\\nuChmLr1drdYXgLmF2bPkDWRgsEECgYEAwf5OrzTSbF3JXZYk2"+
		"RkYouY7/ImcIgPm\\n/6x+zDclQPwuJtT43kCsyVA7SRhGjX3RpLYcFk9Deem0RKAr6i7KJo4s0yAltdzA\\nno089UcjJ068AVLCenOVm2B"+
		"bJV6rywtPJJwmVpQMQrH0Xksif+hPGpxCIzQGXrQQ\\nzdOQpOPm5TsCgYBHgygA8UdBISdy6VGQ+bky6aI0IxGmiIW3N5qrp1PJDhORjx"+
		"nw\\nMr0rYRUYeReZ+2UocDY3m/SVRzYRVLqquVBrCgwUXpgqwIcdcGB4ZgOARZ+xjSKt\\nrwkXJNUS0dVhWA4fbdxALGySgJR29wjAtgxX"+
		"5UfuV6Pwvz5ZB/KDmdJggQKBgDi5\\nPpK2lEzBg67Mx0t/rhd70NCAAFpl37ak3pKiEU+WLXyHS5nZOWzH+/3cjkyzHIjY\\nAxB27tkIAAE"+
		"NAKpCMjPh4LN/M+ege+YgkFF8Eohc2lZct6cMgxNismQT8ZG2Zdbj\\nncY1FfyugjDMMXNLH049oI0gmjg42K0GjsXYKdyfAoGAV1dDKrA07"+
		"CoJ46jtGFN1\\n/clQG5Onm2bXSgwmEETvZyyir8yKKWkHXIojRbU9m8cuHwLbMtF1VeQ0VZcABMZB\\n4/sg+YBXaaHewepJxwei20ewgj4"+
		"SK/togka/kUfyXcKu8kHzOXIeX780EOPMferT\\n32LcrXUFphUzdX6ThYvWxyg=\\n-----END PRIVATE KEY-----"), "\\n", "\n", -1)
	asc.PublicKey = strings.Replace(getStringEnv("PUBLIC_KEY", "-----BEGIN PUBLIC KEY-----\\nMIIBIjANBg" +
		"kqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwKlV08AxrjPvO6X/Ts+t\\n9rWUNo+o4ZEqDRGB2JYU6gTTBD/zf+nasLSc8CggURkwtjfPgVCW" +
		"8p8esVzGgpoy\\n0PrUUmDH3MH9FKGVrKZfLgwgfzmODWcHi14RRLhzE1f7vDOm2YSCCrv/UWgmi8pS\\nYFw14U2fIhTD+WE6xxBV8VaMDF" +
		"RjQmCq5VrXie608nSav9SaOprldp1zlw0puSqz\\nibtrFS6Y1yHnobqSuZ+z63Cvw69UW4rIobyoj7vteGBHa3Xz+6+PirlYWsfTywe2\\nEU" +
		"dkt7bbiXJJAMBLnr8dYT1p967bKYPGYun7vuuS5ZO40xzx3sazGl3fJ+6snVjD\\n+wIDAQAB\\n-----END PUBLIC KEY-----"), "\\n", "\n", -1)
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
