package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	originalEnv := map[string]string{
		"PORT":        os.Getenv("PORT"),
		"LOG_LEVEL":   os.Getenv("LOG_LEVEL"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"AWS_REGION":  os.Getenv("AWS_REGION"),
	}

	// Restore env vars after test
	defer func() {
		for k, v := range originalEnv {
			os.Setenv(k, v)
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name: "default values",
			envVars: map[string]string{},
			expected: &Config{
				Port:        "8080",
				LogLevel:    "info",
				DatabaseURL: "postgres://devuser:devpass@localhost:5432/jobsdb?sslmode=disable",
				SQSEndpoint: "",
				SQSQueueURL: "",
				AWSRegion:   "us-east-2",
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"PORT":         "9090",
				"LOG_LEVEL":    "debug",
				"DB_HOST":      "custom-host",
				"DB_PORT":      "5433",
				"DB_USER":      "testuser",
				"DB_PASSWORD":  "testpass",
				"DB_NAME":      "testdb",
				"SQS_ENDPOINT": "http://localhost:4566",
				"SQS_QUEUE_URL": "http://localhost:4566/queue",
				"AWS_REGION":   "eu-west-1",
			},
			expected: &Config{
				Port:        "9090",
				LogLevel:    "debug",
				DatabaseURL: "postgres://testuser:testpass@custom-host:5433/testdb?sslmode=disable",
				SQSEndpoint: "http://localhost:4566",
				SQSQueueURL: "http://localhost:4566/queue",
				AWSRegion:   "eu-west-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars
			for k := range originalEnv {
				os.Unsetenv(k)
			}

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg := Load()
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestBuildDatabaseURL(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "default values",
			envVars:  map[string]string{},
			expected: "postgres://devuser:devpass@localhost:5432/jobsdb?sslmode=disable",
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "5433",
				"DB_USER":     "admin",
				"DB_PASSWORD": "secret",
				"DB_NAME":     "production",
			},
			expected: "postgres://admin:secret@db.example.com:5433/production?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore env vars
			saved := make(map[string]string)
			for k := range tt.envVars {
				saved[k] = os.Getenv(k)
				os.Setenv(k, tt.envVars[k])
			}
			defer func() {
				for k, v := range saved {
					if v == "" {
						os.Unsetenv(k)
					} else {
						os.Setenv(k, v)
					}
				}
			}()

			result := buildDatabaseURL()
			assert.Equal(t, tt.expected, result)
		})
	}
}