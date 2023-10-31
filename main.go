package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const s3FilePrefix = "mw-code-tester/milan-thakor"

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Couldn't load .env file %v\n", err)
	}

	// Initialize the configuration.
	cfg := NewConfig()

	// Create an AWS session
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(cfg.awsRegion),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		config:     cfg,
		awsSession: sess,
	}

	// Initialize the Gin router.
	r := gin.Default()

	// Simple health check endpoint.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is up and running"})
	})

	// Retrieve top error
	r.GET("/top-error", InjectApp(app, getTopErrorHandler))
	// Upload logs
	r.POST("/upload", InjectApp(app, uploadLogFileHandler))

	// Start the server on the specified port.
	addr := cfg.hostname + ":" + cfg.serverPort
	log.Printf("Server is running on port %s...\n", cfg.serverPort)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
