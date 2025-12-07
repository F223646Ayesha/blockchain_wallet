package database

import (
	"context"
	"crypto-wallet/config"
	"crypto-wallet/models"
	"time"
)

// Save a log entry
func AddLog(level, message string) {
	ctx := context.Background()

	entry := models.LogEntry{
		Timestamp: time.Now().Unix(),
		Level:     level,
		Message:   message,
	}

	config.Firestore.Collection("logs").Add(ctx, entry)
}

// Read logs
func GetSystemLogs() ([]models.LogEntry, error) {

	ctx := context.Background()

	snaps, err := config.Firestore.Collection("logs").
		OrderBy("timestamp", 1).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, err
	}

	var logs []models.LogEntry
	for _, d := range snaps {
		var entry models.LogEntry
		d.DataTo(&entry)
		logs = append(logs, entry)
	}

	return logs, nil
}
