package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port         string
	LogLevel     string
	DatabaseURL  string
	SQSEndpoint  string
	SQSQueueURL  string
	AWSRegion    string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DatabaseURL: buildDatabaseURL(),
		SQSEndpoint: getEnv("SQS_ENDPOINT", ""),
		SQSQueueURL: getEnv("SQS_QUEUE_URL", ""),
		AWSRegion:   getEnv("AWS_REGION", "us-east-2"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildDatabaseURL() string {
	// Check for DATABASE_URL first
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	
	// Fall back to building from individual components
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "devuser")
	pass := getEnv("DB_PASSWORD", "devpass")
	name := getEnv("DB_NAME", "jobsdb")
	
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
}