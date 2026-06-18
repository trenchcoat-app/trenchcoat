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

var AppConfig Config

func Init() error {
	err := godotenv.Load("../.env")
	if err != nil {
		return err
	}

	AppConfig = Config{
		POSTGRES_USER:        getEnvOrPanic("POSTGRES_USER"),
		POSTGRES_PASSWORD:    getEnvOrPanic("POSTGRES_PASSWORD"),
		POSTGRES_DB:          getEnvOrPanic("POSTGRES_DB"),
		POSTGRES_HOST:        getEnvOrDefault("POSTGRES_HOST", "localhost"),
		POSTGRES_PORT:        getEnvOrDefault("POSTGRES_PORT", "5432"),
		CORS_ALLOWED_ORIGINS: strings.Split(getEnvOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ","),
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
