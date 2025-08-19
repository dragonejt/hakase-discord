// Package clients implements backend API operations for courses.
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

type Course struct {
	ID            int    `json:"id,omitempty"`
	Platform      int    `json:"platform,omitempty"`
	CourseID      string `json:"course_id"`
	NotifyChannel string `json:"notify_channel,omitempty"`
	NotifyGroup   string `json:"notify_group,omitempty"`
}

// ReadCourse retrieves a course by its ID from the backend.
func (backend *BackendClient) ReadCourse(span *sentry.Span, courseID string) (Course, error) {
	span = span.StartChild("readCourse")
	defer span.Finish()

	course := Course{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/courses?course_id=%s", backend.URL, courseID), nil)
	if err != nil {
		return course, stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return course, stacktrace.Propagate(err, "failed to execute API request")
	}
	if response.StatusCode != http.StatusOK {
		return course, stacktrace.Propagate(err, "failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return course, stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &course)
	if err != nil {
		return course, stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}

	return course, nil
}

// HeadCourse checks if a course exists in the backend.
func (backend *BackendClient) HeadCourse(span *sentry.Span, courseID string) error {
	span = span.StartChild("headCourse")
	defer span.Finish()

	request, err := http.NewRequest(http.MethodHead, fmt.Sprintf("%s/courses?course_id=%s", backend.URL, courseID), nil)
	if err != nil {
		return stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return stacktrace.Propagate(err, "failed to execute API request")
	}
	if response.StatusCode != http.StatusOK {
		return stacktrace.Propagate(err, "failed status code API response: %d", response.StatusCode)
	}

	return nil
}

// CreateCourse creates a new course in the backend.
func (backend *BackendClient) CreateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("createCourse")
	defer span.Finish()

	jsonBody, err := json.Marshal(course)
	if err != nil {
		return stacktrace.Propagate(err, "failed to marshal course")
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/courses", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return stacktrace.Propagate(err, "failed to execute API request")
	}
	if response.StatusCode != http.StatusCreated {
		return stacktrace.Propagate(err, "failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &course)
	if err != nil {
		return stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}

	return nil
}

// UpdateCourse updates an existing course in the backend.
func (backend *BackendClient) UpdateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("updateCourse")
	defer span.Finish()

	jsonBody, err := json.Marshal(course)
	if err != nil {
		return stacktrace.Propagate(err, "failed to marshal course")
	}

	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/courses", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return stacktrace.Propagate(err, "failed to execute API request")
	}
	if response.StatusCode != http.StatusAccepted {
		return stacktrace.Propagate(err, "failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return stacktrace.Propagate(err, "failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &course)
	if err != nil {
		return stacktrace.Propagate(err, "failed to unmarshal API response: %s", string(body))
	}

	return nil
}

// DeleteCourse deletes a course from the backend.
func (backend *BackendClient) DeleteCourse(span *sentry.Span, courseID string) error {
	span = span.StartChild("deleteCourse")
	defer span.Finish()

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/courses?course_id=%s", backend.URL, courseID), nil)
	if err != nil {
		return stacktrace.Propagate(err, "failed to create API request")
	}
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return stacktrace.Propagate(err, "failed to execute API request")
	}
	if response.StatusCode != http.StatusNoContent {
		return stacktrace.Propagate(err, "failed status code API response: %d", response.StatusCode)
	}

	return nil
}
