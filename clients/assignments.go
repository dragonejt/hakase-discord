package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

type Assignment struct {
	ID       int       `json:"id,omitempty"`
	Course   int       `json:"course,omitempty"`
	CourseID string    `json:"course_id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Due      time.Time `json:"due,omitempty"`
	Link     string    `json:"link,omitempty"`
}

func (backend *BackendClient) ReadAssignment(span *sentry.Span, assignmentID string) (Assignment, error) {
	span = span.StartChild("readAssignment")
	defer span.Finish()

	assignment := Assignment{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/assignments?id=%s", backend.URL, assignmentID), nil)
	if err != nil {
		return assignment, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return assignment, fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return assignment, fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return assignment, fmt.Errorf("failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &assignment)
	if err != nil {
		return assignment, fmt.Errorf("failed to unmarshal API response: %s", string(body))
	}

	return assignment, nil
}

func (backend *BackendClient) HeadAssignment(span *sentry.Span, assignmentID string) error {
	span = span.StartChild("headAssignment")
	defer span.Finish()

	request, err := http.NewRequest(http.MethodHead, fmt.Sprintf("%s/assignments?id=%s", backend.URL, assignmentID), nil)
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	return nil
}

func (backend *BackendClient) ListAssignments(span *sentry.Span, courseID string) ([]Assignment, error) {
	span = span.StartChild("listAssignments")
	defer span.Finish()

	assignments := []Assignment{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/assignments?course_id=%s", backend.URL, courseID), nil)
	if err != nil {
		return assignments, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return assignments, fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return assignments, fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return assignments, fmt.Errorf("failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &assignments)
	if err != nil {
		return assignments, fmt.Errorf("failed to unmarshal API response: %s", string(body))
	}

	return assignments, nil
}

func (backend *BackendClient) CreateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {
	span = span.StartChild("createAssignment")
	defer span.Finish()

	jsonBody, err := json.Marshal(assignment)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to marshal assignment: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/assignments", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusCreated {
		return Assignment{}, fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &assignment)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to unmarshal API response: %s", string(body))
	}

	return assignment, nil
}

func (backend *BackendClient) UpdateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {
	span = span.StartChild("updateAssignment")
	defer span.Finish()

	jsonBody, err := json.Marshal(assignment)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to marshal assignment: %w", err)
	}

	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/assignments", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusAccepted {
		return Assignment{}, fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &assignment)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to unmarshal API response: %s", string(body))
	}

	return assignment, nil
}

func (backend *BackendClient) DeleteAssignment(span *sentry.Span, assignmentID string) error {
	span = span.StartChild("deleteAssignment")
	defer span.Finish()

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/assignments?id=%s", backend.URL, assignmentID), nil)
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	return nil
}
