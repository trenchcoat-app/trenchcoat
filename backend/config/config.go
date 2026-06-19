package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_USER        string
	POSTGRES_PASSWORD    string
	POSTGRES_DB          string
	POSTGRES_HOST        string
	POSTGRES_PORT        string
	CORS_ALLOWED_ORIGINS []string
}

var defaultConfig = Config{
	POSTGRES_HOST:        "localhost",
	POSTGRES_PORT:        "5432",
	CORS_ALLOWED_ORIGINS: []string{"http://localhost:5173"},
}
var AppConfig = defaultConfig

func Init() error {
	err := godotenv.Load("../.env")
	if err != nil {
		return err
	}

	AppConfig = Config{
		POSTGRES_USER:        getEnvOrPanic("POSTGRES_USER"),
		POSTGRES_PASSWORD:    getEnvOrPanic("POSTGRES_PASSWORD"),
		POSTGRES_DB:          getEnvOrPanic("POSTGRES_DB"),
		POSTGRES_HOST:        getEnvOrDefault("POSTGRES_HOST", defaultConfig.POSTGRES_HOST),
		POSTGRES_PORT:        getEnvOrDefault("POSTGRES_PORT", defaultConfig.POSTGRES_PORT),
		CORS_ALLOWED_ORIGINS: getEnvOrDefaultSlice("CORS_ALLOWED_ORIGINS", defaultConfig.CORS_ALLOWED_ORIGINS),
	}

	return nil
}

func getEnvOrPanic(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("required env var " + key + " is not set")
	}
	return val
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvOrDefaultSlice(key string, fallback []string) []string {
	if val := os.Getenv(key); val != "" {
		return strings.Split(val, ",")
	}
	return fallback
}
