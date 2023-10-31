package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// topError represents the top error structure with relevant details.
type topError struct {
	Message  string `json:"top-error"`
	Count    int    `json:"count"`
	Service  string `json:"service"`
	Severity string `json:"severity"`
}

// extractTopErrorFromS3Bucket extracts the top error from S3 bucket logs.
func extractTopErrorFromS3Bucket(svc *s3.S3, bucketName string, result *s3.ListObjectsOutput) (topError, error) {
	var topErr topError
	for _, item := range result.Contents {
		fileDirs := strings.Split(*item.Key, "/")
		if len(fileDirs) < 6 {
			return topError{}, errors.New("invalid s3 file path")
		}
		service := fileDirs[3]
		severity := fileDirs[4]

		fileContents, err := getFileContentsFromS3(svc, bucketName, *item.Key)
		if err != nil {
			return topError{}, err
		}

		newTopErr, err := getTopError(fileContents, service, severity, topErr)
		if err != nil {
			return topError{}, err
		}

		topErr = newTopErr
	}

	return topErr, nil
}

// getFileContentsFromS3 retrieves file contents from the S3 bucket.
func getFileContentsFromS3(svc *s3.S3, bucketName, filepath string) ([]string, error) {
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filepath),
	})
	if err != nil {
		return nil, errors.New("failed to load s3 bucket file")
	}
	defer result.Body.Close()

	var fileContents []string
	scanner := bufio.NewScanner(result.Body)
	for scanner.Scan() {
		fileContents = append(fileContents, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.New("failed to read s3 file")
	}

	return fileContents, nil
}

// getTopError identifies the top error from the log contents.
func getTopError(fileContents []string, service, severity string, currentTopErr topError) (topError, error) {
	for _, line := range fileContents {
		lineParts := strings.SplitN(line, " - ", 2)
		if len(lineParts) < 2 {
			return topError{}, errors.New("invalid log format")
		}

		logCount, err := strconv.Atoi(lineParts[0])
		if err != nil {
			return topError{}, errors.New("invalid log count")
		}

		if logCount > currentTopErr.Count {
			currentTopErr = topError{
				Message:  lineParts[1],
				Count:    logCount,
				Service:  service,
				Severity: severity,
			}
		}
	}

	return currentTopErr, nil
}

// uploadLogsToS3 uploads logs to the specified S3 bucket.
func uploadLogsToS3(sess *session.Session, bucketName string, logsFiles map[string]string) error {
	uploader := s3manager.NewUploader(sess)

	for filepath, logs := range logsFiles {
		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(filepath),
			Body:   strings.NewReader(logs),
		})
		if err != nil {
			return fmt.Errorf("failed to upload file, %v", err)
		}
	}

	return nil
}
