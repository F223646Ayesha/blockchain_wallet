package services

import (
	"context"
	"crypto-wallet/config"
	"crypto-wallet/models"
)

// GetLogsService returns all logs stored for debugging, mining, errors, etc.
func GetLogsService() ([]models.LogEntry, error) {
	ctx := context.Background()

	// Query "logs" collection
	snaps, err := config.Firestore.Collection("logs").OrderBy("timestamp", 1).Documents(ctx).GetAll()
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
