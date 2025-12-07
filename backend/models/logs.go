package models

type LogEntry struct {
	Timestamp int64  `json:"timestamp" firestore:"timestamp"`
	Level     string `json:"level" firestore:"level"`
	Message   string `json:"message" firestore:"message"`
}
