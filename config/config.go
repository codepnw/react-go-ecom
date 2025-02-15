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
	*JWTConfig
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

type JWTConfig struct {
	Secret             string
	AccessTokenExpire  int
	RefreshTokenExpire int
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
		&JWTConfig{
			Secret:             getEnv("JWT_SECRET", "secret"),
			AccessTokenExpire:  getEnvInt("JWT_ACCESS_EXPIRE", 15),
			RefreshTokenExpire: getEnvInt("JWT_REFRESH_EXPIRE", 1440),
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
