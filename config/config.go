package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	*AppConfig
	*DBConfig
}

type AppConfig struct {
	AppPort    string
	AppVersion string
}

type DBConfig struct {
	DBAddr       string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func LoadConfig(envPath string) *Config {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("cant loading .env file:", err)
	}

	return &Config{
		&AppConfig{
			AppPort:    getEnv("APP_PORT", "8080"),
			AppVersion: getEnv("APP_VERSION", "v1"),
		},
		&DBConfig{
			DBAddr:       getEnv("DB_URL", ""),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxIdleTime:  getEnv("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if v, exists := os.LookupEnv(key); exists {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if v, exists := os.LookupEnv(key); exists {
		value, _ := strconv.Atoi(v)
		return value
	}
	return defaultValue
}
