package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type Course struct {
	ID            int    `json:"id,omitempty"`
	Platform      int    `json:"platform,omitempty"`
	CourseID      string `json:"course_id"`
	NotifyChannel string `json:"notify_channel,omitempty"`
	NotifyGroup   string `json:"notify_group,omitempty"`
}

func (backend *BackendClient) ReadCourse(span *sentry.Span, courseID string) (Course, error) {
	span = span.StartChild("readCourse")
	defer span.Finish()

	course := Course{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/courses?course_id=%s", backend.URL, courseID), nil)
	if err != nil {
		return course, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return course, fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return course, fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return course, fmt.Errorf("failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &course)
	if err != nil {
		return course, fmt.Errorf("failed to unmarshal API response: %s", string(body))
	}

	return course, nil
}

func (backend *BackendClient) HeadCourse(span *sentry.Span, courseID string) error {
	span = span.StartChild("headCourse")
	defer span.Finish()

	request, err := http.NewRequest(http.MethodHead, fmt.Sprintf("%s/courses?course_id=%s", backend.URL, courseID), nil)
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

func (backend *BackendClient) CreateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("createCourse")
	defer span.Finish()

	jsonBody, err := json.Marshal(course)
	if err != nil {
		return fmt.Errorf("failed to marshal course: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/courses", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &course)
	if err != nil {
		return fmt.Errorf("failed to unmarshal API response: %s", string(body))
	}

	return nil
}

func (backend *BackendClient) UpdateCourse(span *sentry.Span, course Course) error {
	span = span.StartChild("updateCourse")
	defer span.Finish()

	jsonBody, err := json.Marshal(course)
	if err != nil {
		return fmt.Errorf("failed to marshal course: %w", err)
	}

	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/courses", backend.URL), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", fmt.Sprintf("Token %s", backend.API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	response, err := backend.HTTP_CLIENT.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &course)
	if err != nil {
		return fmt.Errorf("failed to unmarshal API response: %s", string(body))
	}

	return nil
}

func (backend *BackendClient) DeleteCourse(span *sentry.Span, courseID string) error {
	span = span.StartChild("deleteCourse")
	defer span.Finish()

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/courses?course_id=%s", backend.URL, courseID), nil)
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
