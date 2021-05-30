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

	/*asc.PrivateKey = strings.Replace(getStringEnv("PRIVATE_KEY", "-----BEGIN PRIVATE KEY-----\\nMIICdwIBADANBgkqhkiG"+
	"9w0BAQEFAASCAmEwggJdAgEAAoGBANnmLifeLBsiXe/J\\n8O3ophHHaCfJ+EdAUYn7vArJTUtankCD3I8O3n+QM0KNsXzXd+eN6VmNm3bjLp"+
	"Hq\\nVjI/jCr2m1EqXgvRQP74/wOU1sHN3zSRQbcPR0dfJiDfTRmfh/LVrKgcU0kQ4yrG\\nlc0KGB2uslzrKLJCmQ4G0WeM3tKNAgMBAAECgY"+
	"EApsep+FXzSGmLoOfegxqZUe5g\\n6GOMp2yxfH2ztkXR5aVcj2DeRplI8DZ9Jamyei2p1xAl1aevoNXOZV0J0LgXHbm+\\nP6MGU7d+IYD2hI"+
	"CWPfD4pqJafkYc7Q94eQaIiShlYEOoEiLDt09m2V3J/VWxEWw0\\nGTzT1T6zDuwD5epXY/kCQQDv+Xeq+SU5+avfysvm/8bITu/WBRXKxQ7V"+
	"2dg9rJIF\\nrAZSTPUIqdKm2F+o8DIX4sSMouFMgo81Ad4S8D13iQCPAkEA6HNRwmcfQCHzsuuT\\n1407mEPFAcgIckU6e9ubXRRepWPjE6MJ"+
	"IyrDeIkCJfgFPiK8OcNvFLUCD8NaySD1\\nQuHRIwJAJtIqo8QOW6SiQ1/hQItcMwdiETNdZSIf1kSZkNCcBsLfeuzsLuyaIVeb\nkg7Za7fJ"+
	"qB6pZ+EvHZohvNqUdwP4zQJBALf6piiG9C4PcVIYsOA3cYa3hNM/HqhK\\n8NodW9+VAsBGyfC95rqF20aosiGZJ5UhavcRHvc1uNb/GPj99A"+
	"EmuB8CQCy9M89/\\nWGs0V60TrOWn2cmNlvexvxJgtWIjzdtp5rBj/E7Dmfx9nE6sG+uJqob389HYb0fF\\nj8MrN6RCirNhupc=\\n-----END"+
	" PRIVATE KEY-----"), "\\n", "\n", -1)*/
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
