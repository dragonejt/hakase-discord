package clients

import (
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

type HakaseClient interface {
	// Course APIs
	ReadCourse(span *sentry.Span, courseID string) (Course, error)
	CreateCourse(span *sentry.Span, course Course) error
	UpdateCourse(span *sentry.Span, course Course) error
	DeleteCourse(span *sentry.Span, courseID string) error
	// Assignment APIs
	ReadAssignment(span *sentry.Span, assignmentID string) (Assignment, error)
	ListAssignments(span *sentry.Span, courseID string) ([]Assignment, error)
	CreateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error)
	UpdateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error)
	DeleteAssignment(span *sentry.Span, assignmentID string) error
}

type BackendClient struct {
	HakaseClient
	URL         string
	API_KEY     string
	HTTP_CLIENT *http.Client
}

type DiscordSession struct{}

type AssignmentNotification struct {
	AssignmentID int
	Before       time.Duration
}

type StudySessionNotification struct {
	SessionID int
	Timestamp time.Time
}
