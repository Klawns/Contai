package app

import (
	"fmt"
	"os"
)

type config struct {
	port        string
	databaseDSN string
	jwtSecret   string
}

func loadConfig() (config, error) {
	jwtSecret, err := getJWTSecret()
	if err != nil {
		return config{}, err
	}

	return config{
		port:        getEnv("PORT", "8080"),
		databaseDSN: getDatabaseDSN(),
		jwtSecret:   jwtSecret,
	}, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func isProduction() bool {
	return os.Getenv("APP_ENV") == "production" || os.Getenv("GIN_MODE") == "release"
}

func getJWTSecret() (string, error) {
	value := os.Getenv("JWT_SECRET")
	if value == "" && isProduction() {
		return "", fmt.Errorf("JWT_SECRET is required in production")
	}
	if value == "" {
		return "dev-secret-change-me", nil
	}

	return value, nil
}

func getDatabaseDSN() string {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		return value
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "contai"),
		getEnv("DB_PASSWORD", "contai"),
		getEnv("DB_NAME", "contai"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
		getEnv("DB_TIMEZONE", "America/Sao_Paulo"),
	)
}
