package notifications

import "time"

type AssignmentNotification struct {
	AssignmentID int
	Before       time.Duration
}

type StudySessionNotification struct {
	SessionID int
	Timestamp time.Time
}
