package config

import (
	"auth-service/pkg/logging"
	"go.uber.org/zap"
	"os"
	"strconv"
)

var logger *zap.Logger

func init() {
	logger = logging.GetLogger()
}

// GetStringEnv gets the specific environment variables with default value, returns default value if variable not set
func GetStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// GetIntEnv gets the specific environment variables with default value, returns default value if variable not set
func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return convertStringToInt(value)
}

func convertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Warn("an error occured while converting from string to int. Setting it as zero",
			zap.String("error", err.Error()))
		i = 0
	}
	return i
}
