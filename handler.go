package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// getTopErrorHandler handles the retrieval of top error from the S3.
func getTopErrorHandler(ctx *gin.Context, app *application) {
	svc := s3.New(app.awsSession)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(app.config.awsBucketName),
		Prefix: aws.String(s3FilePrefix),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		log.Printf("failed to load s3 bucket: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load s3 bucket"})
		return
	}

	topErr, err := extractTopErrorFromS3Bucket(svc, app.config.awsBucketName, result)
	if err != nil {
		log.Printf("failed to extract top-error from s3 bucket: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to extract top-error from s3 bucket"})
		return
	}

	ctx.JSON(http.StatusOK, topErr)
}

// uploadLogFileHandler handles the uploading of log files to S3.
func uploadLogFileHandler(ctx *gin.Context, app *application) {
	file, _, err := ctx.Request.FormFile("sample-file")
	if err != nil {
		log.Printf("failed to load file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	logs, err := processLogFile(file)
	if err != nil {
		log.Printf("failed to process log file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process log file"})
		return
	}

	s3LogsFiles := parseLogsForS3(logs)

	if err := uploadLogsToS3(app.awsSession, app.config.awsBucketName, s3LogsFiles); err != nil {
		log.Printf("failed to upload logs to s3: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload logs to s3"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "logs uploaded to s3 successfully",
	})
}
