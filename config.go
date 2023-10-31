package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
)

// Application holds application-wide dependencies.
type application struct {
	config     *config
	awsSession *session.Session
}

// Config holds the application configuration.
type config struct {
	awsRegion     string
	awsBucketName string
	hostname      string
	serverPort    string
}

// New creates a new Config instance with the default values.
func NewConfig() *config {
	return &config{
		awsRegion:     getEnv("AWS_REGION", "us-east-1"),
		awsBucketName: getEnv("AWS_BUCKET_NAME", "mw-code-tester"),
		hostname:      getEnv("HOSTNAME", "localhost"),
		serverPort:    getEnv("SERVER_PORT", "8080"),
	}
}

// getEnv retrieves an environment variable with a default value if not set.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}
