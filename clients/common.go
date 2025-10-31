// Package clients provides interfaces and types for interacting with the backend API and Discord session.
package clients

import (
	"net/http"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

type HakaseClient struct {
	Backend       BackendClient
	Notifications NotificationsClient
}

type BackendClient interface {
	// Course APIs
	ReadCourse(span *sentry.Span, courseID string) (Course, error)
	HeadCourse(span *sentry.Span, courseID string) error
	CreateCourse(span *sentry.Span, course Course) error
	UpdateCourse(span *sentry.Span, course Course) error
	DeleteCourse(span *sentry.Span, courseID string) error
	// Assignment APIs
	ReadAssignment(span *sentry.Span, assignmentID string) (Assignment, error)
	HeadAssignment(span *sentry.Span, assignmentID string) error
	ListAssignments(span *sentry.Span, courseID string) ([]Assignment, error)
	CreateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error)
	UpdateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error)
	DeleteAssignment(span *sentry.Span, assignmentID string) error
}

type NotificationsClient interface {
	ListenToStream(bot *discordgo.Session, hakaseClient BackendClient, stopListener chan bool)
	PublishNotification(span *sentry.Span, notification string)
	PublishAssignmentNotification(span *sentry.Span, notification AssignmentNotification)
	PublishStudySessionNotification(span *sentry.Span, notification StudySessionNotification)
}

type APIClient struct {
	BackendClient
	Url        string
	APIKey     string
	HttpClient *http.Client
}

type MQClient struct {
	NotificationsClient
	NATSUrl       string
	StreamName    string
	PublisherPool sync.Pool
}

type DiscordSession struct{}

type AssignmentNotification struct {
	AssignmentID int
	CourseID     string
	Before       time.Duration
}

type StudySessionNotification struct {
	SessionID int
	CourseID  string
	Timestamp time.Time
}
