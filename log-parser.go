package main

import (
	"bufio"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"
)

const (
	timestampFormat = "2006-01-02T15:04:05Z"
)

// logMessage represents the log message structure.
type logMessage struct {
	timestamp time.Time
	severity  string
	service   string
	message   string
}

// parseLogEntry parses a single log entry from a string line.
func parseLogEntry(line string) (logMessage, error) {
	lineParts := strings.SplitN(line, " ", 4)
	if len(lineParts) < 4 {
		return logMessage{}, errors.New("invalid log format")
	}

	ts, err := time.Parse(timestampFormat, lineParts[0])
	if err != nil {
		return logMessage{}, errors.New("invalid timestamp format")
	}

	return logMessage{
		timestamp: ts,
		severity:  lineParts[1],
		service:   strings.Trim(lineParts[2], "[]"),
		message:   lineParts[3],
	}, nil
}

// processLogFile reads a multipart file, processes log entries, and returns a slice of logMessages.
func processLogFile(file multipart.File) ([]logMessage, error) {
	var logs []logMessage

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logEntry, err := parseLogEntry(scanner.Text())
		if err != nil {
			return nil, err
		}
		logs = append(logs, logEntry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// parseLogsForS3 categorizes log messages by time window, service, severity, and creates a map for S3 file storage.
func parseLogsForS3(logs []logMessage) map[string]string {
	logsByCategory := make(map[string]map[string]int)

	// Categorize logs by time window, service, severity, and count
	for _, log := range logs {
		logWindowStart := log.timestamp.Truncate(time.Hour)
		logWindowEnd := logWindowStart.Add(time.Hour)
		logWindow := logWindowStart.Format("2006-01-02T15:04:05") + "-" + logWindowEnd.Format("15:04:05")

		category := fmt.Sprintf("%s/%s/%s/%s/sample.log", s3FilePrefix, logWindow, log.service, log.severity)

		if _, ok := logsByCategory[category]; !ok {
			logsByCategory[category] = map[string]int{
				log.message: 1,
			}
		} else {
			logsByCategory[category][log.message]++
		}
	}

	s3LogsFiles := make(map[string]string)
	for filename, logs := range logsByCategory {
		var messages []string
		for msg, count := range logs {
			messages = append(messages, fmt.Sprintf("%d - %s", count, msg))
		}

		s3LogsFiles[filename] = strings.Join(messages, "\n")
	}

	return s3LogsFiles
}
