package clients

import "time"

type DiscordSession struct{}

type AssignmentNotification struct {
	AssignmentID int
	Before       time.Duration
}

type StudySessionNotification struct {
	SessionID int
	Timestamp time.Time
}
