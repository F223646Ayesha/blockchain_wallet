package utils

import (
	"crypto-wallet/config"
	"time"
)

func LogEvent(action string, note string) {
	config.Firestore.Collection("logs").Add(nil, map[string]interface{}{
		"action":    action,
		"note":      note,
		"timestamp": time.Now().Unix(),
	})
}
