package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_USER          string
	POSTGRES_PASSWORD      string
	POSTGRES_DB            string
	POSTGRES_HOST          string
	POSTGRES_PORT          string
	CORS_ALLOWED_ORIGINS   []string
	COOKIE_NAME            string
	COOKIE_DOMAIN          string
	COOKIE_SECURE          bool
	COOKIE_SAME_SITE       string
	COOKIE_PATH            string
	SESSION_EXPIRY_SECONDS int
}

var defaultConfig = Config{
	POSTGRES_HOST:          "localhost",
	POSTGRES_PORT:          "5432",
	CORS_ALLOWED_ORIGINS:   []string{"http://localhost:5173"},
	COOKIE_NAME:            "sid",
	COOKIE_DOMAIN:          "",
	COOKIE_SECURE:          false,
	COOKIE_SAME_SITE:       "Lax",
	COOKIE_PATH:            "/",
	SESSION_EXPIRY_SECONDS: 60 * 60 * 24,
}
var AppConfig = defaultConfig

func Init() error {
	err := godotenv.Load("../.env")
	if err != nil {
		return err
	}

	AppConfig = Config{
		POSTGRES_USER:          getEnvOrPanic("POSTGRES_USER"),
		POSTGRES_PASSWORD:      getEnvOrPanic("POSTGRES_PASSWORD"),
		POSTGRES_DB:            getEnvOrPanic("POSTGRES_DB"),
		POSTGRES_HOST:          getEnvOrDefaultString("POSTGRES_HOST", defaultConfig.POSTGRES_HOST),
		POSTGRES_PORT:          getEnvOrDefaultString("POSTGRES_PORT", defaultConfig.POSTGRES_PORT),
		CORS_ALLOWED_ORIGINS:   getEnvOrDefaultSlice("CORS_ALLOWED_ORIGINS", defaultConfig.CORS_ALLOWED_ORIGINS),
		COOKIE_NAME:            getEnvOrDefaultString("COOKIE_NAME", defaultConfig.COOKIE_NAME),
		COOKIE_DOMAIN:          getEnvOrDefaultString("COOKIE_DOMAIN", defaultConfig.COOKIE_DOMAIN),
		COOKIE_SECURE:          getEnvOrDefaultBool("COOKIE_SECURE", defaultConfig.COOKIE_SECURE),
		COOKIE_SAME_SITE:       getEnvOrDefaultString("COOKIE_SAME_SITE", defaultConfig.COOKIE_SAME_SITE),
		COOKIE_PATH:            getEnvOrDefaultString("COOKIE_PATH", defaultConfig.COOKIE_PATH),
		SESSION_EXPIRY_SECONDS: getEnvOrDefaultInt("SESSION_EXPIRY_SECONDS", defaultConfig.SESSION_EXPIRY_SECONDS),
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

func getEnvOrDefaultString(key string, fallback string) string {
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

// Gets env var as bool if it exists, defaults to `fallback` if it doesn't
// This can still panic if the env var value doesn't parse to bool
func getEnvOrDefaultBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	switch val {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		panic("Failed to parse " + key + " into bool: unexpected value " + val)
	}
}

// Gets env var as int if it exists, defaults to `fallback` if it doesn't
// This can still panic if the env var value doesn't parse to int
func getEnvOrDefaultInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		parsedInt, err := strconv.ParseInt(val, 10, 0)

		if err != nil {
			panic("Failed to parse " + key + " into int: " + err.Error())
		}

		// No truncation can happen due to bitSize=0 in the ParseInt call
		return int(parsedInt)
	}
	return fallback
}
