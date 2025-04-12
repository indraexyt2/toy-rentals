package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	// Server
	ServerPort string
	IsProd     bool

	// Database
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func LoadConfig() *Config {
	// Load .env file jika ada
	godotenv.Load()

	return &Config{
		// Server
		ServerPort: getEnv("SERVER_PORT", "8080"),
		IsProd:     getEnvAsBool("IS_PROD", false),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "toyrentals"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
