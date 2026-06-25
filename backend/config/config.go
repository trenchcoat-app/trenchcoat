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
	COOKIE_NAME          string
	COOKIE_DOMAIN        string
	COOKIE_SECURE        bool
	COOKIE_SAME_SITE     string
	COOKIE_PATH          string
}

var defaultConfig = Config{
	POSTGRES_HOST:        "localhost",
	POSTGRES_PORT:        "5432",
	CORS_ALLOWED_ORIGINS: []string{"http://localhost:5173"},
	COOKIE_NAME:          "sid",
	COOKIE_DOMAIN:        "",
	COOKIE_SECURE:        false,
	COOKIE_SAME_SITE:     "Lax",
	COOKIE_PATH:          "/",
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
		COOKIE_NAME:          getEnvOrDefault("COOKIE_NAME", defaultConfig.COOKIE_NAME),
		COOKIE_DOMAIN:        getEnvOrDefault("COOKIE_DOMAIN", defaultConfig.COOKIE_DOMAIN),
		COOKIE_SECURE:        getEnvOrDefaultBool("COOKIE_SECURE", defaultConfig.COOKIE_SECURE),
		COOKIE_SAME_SITE:     getEnvOrDefault("COOKIE_SAME_SITE", defaultConfig.COOKIE_SAME_SITE),
		COOKIE_PATH:          getEnvOrDefault("COOKIE_PATH", defaultConfig.COOKIE_PATH),
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

func getEnvOrDefaultBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		return val == "true" || val == "1" || val == "yes"
	}
	return fallback
}
