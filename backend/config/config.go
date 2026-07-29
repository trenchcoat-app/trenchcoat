package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_USER            string
	POSTGRES_PASSWORD        string
	POSTGRES_DB              string
	POSTGRES_HOST            string
	POSTGRES_PORT            string
	CORS_ALLOWED_ORIGINS     []string
	SESSION_COOKIE_NAME      string
	SESSION_COOKIE_DOMAIN    string
	SESSION_COOKIE_SECURE    bool
	SESSION_COOKIE_SAME_SITE string
	SESSION_COOKIE_PATH      string
	SESSION_EXPIRY_SECONDS   int
}

/**
 * Default configuration is made to seamlessly work with local dev environments.
 * For production environments, the following **must** be changed:
 *  - SESSION_COOKIE_SECURE: Should be true to require HTTPS
 *  - SESSION_COOKIE_DOMAIN: Should be the frontend domain in prod
 */
var defaultConfig = Config{
	POSTGRES_HOST:            "localhost",
	POSTGRES_PORT:            "5432",
	CORS_ALLOWED_ORIGINS:     []string{"http://localhost:5173"},
	SESSION_COOKIE_NAME:      "sid",
	SESSION_COOKIE_DOMAIN:    "",
	SESSION_COOKIE_SECURE:    false,
	SESSION_COOKIE_SAME_SITE: "Lax",
	SESSION_COOKIE_PATH:      "/",
	SESSION_EXPIRY_SECONDS:   60 * 60 * 24,
}
var AppConfig = defaultConfig

func Init() error {
	err := godotenv.Load("../.env")
	if err != nil {
		return err
	}

	AppConfig = Config{
		POSTGRES_USER:            getEnvOrPanic("POSTGRES_USER"),
		POSTGRES_PASSWORD:        getEnvOrPanic("POSTGRES_PASSWORD"),
		POSTGRES_DB:              getEnvOrPanic("POSTGRES_DB"),
		POSTGRES_HOST:            getEnvOrDefaultString("POSTGRES_HOST", defaultConfig.POSTGRES_HOST),
		POSTGRES_PORT:            getEnvOrDefaultString("POSTGRES_PORT", defaultConfig.POSTGRES_PORT),
		CORS_ALLOWED_ORIGINS:     getEnvOrDefaultSlice("CORS_ALLOWED_ORIGINS", defaultConfig.CORS_ALLOWED_ORIGINS),
		SESSION_COOKIE_NAME:      getEnvOrDefaultString("SESSION_COOKIE_NAME", defaultConfig.SESSION_COOKIE_NAME),
		SESSION_COOKIE_DOMAIN:    getEnvOrDefaultString("SESSION_COOKIE_DOMAIN", defaultConfig.SESSION_COOKIE_DOMAIN),
		SESSION_COOKIE_SECURE:    getEnvOrDefaultBool("SESSION_COOKIE_SECURE", defaultConfig.SESSION_COOKIE_SECURE),
		SESSION_COOKIE_SAME_SITE: getEnvOrDefaultString("SESSION_COOKIE_SAME_SITE", defaultConfig.SESSION_COOKIE_SAME_SITE),
		SESSION_COOKIE_PATH:      getEnvOrDefaultString("SESSION_COOKIE_PATH", defaultConfig.SESSION_COOKIE_PATH),
		SESSION_EXPIRY_SECONDS:   getEnvOrDefaultInt("SESSION_EXPIRY_SECONDS", defaultConfig.SESSION_EXPIRY_SECONDS),
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
