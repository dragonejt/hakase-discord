package notifications

import "time"

type AssignmentNotification struct {
	AssignmentID int
	Timestamp    time.Time
}

type StudySessionNotification struct {
	SessionID int
	Timestamp time.Time
}
