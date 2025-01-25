package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dragonejt/hakase-discord/settings"
	"github.com/getsentry/sentry-go"
)

type Course struct {
	ID            int    `json:"id,omitempty"`
	Platform      int    `json:"platform,omitempty"`
	CourseID      string `json:"course_id"`
	NotifyChannel string `json:"notify_channel,omitempty"`
	NotifyGroup   string `json:"notify_group,omitempty"`
}

func ReadCourse(courseID string) (Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sentry.StartSpan(ctx, "readCourse")
	course := Course{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/courses?course_id=%s", settings.BACKEND_URL, courseID), nil)
	if err != nil {
		return course, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))

	response, err := http.DefaultClient.Do(request)
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

func CreateCourse(course Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sentry.StartSpan(ctx, "createCourse")
	jsonBody, err := json.Marshal(course)
	if err != nil {
		return fmt.Errorf("failed to marshal course: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/courses", settings.BACKEND_URL), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))

	response, err := http.DefaultClient.Do(request)
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

func UpdateCourse(course Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sentry.StartSpan(ctx, "updateCourse")
	jsonBody, err := json.Marshal(course)
	if err != nil {
		return fmt.Errorf("failed to marshal course: %w", err)
	}

	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/courses", settings.BACKEND_URL), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))

	response, err := http.DefaultClient.Do(request)
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

func DeleteCourse(courseID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sentry.StartSpan(ctx, "deleteCourse")
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/courses?course_id=%s", settings.BACKEND_URL, courseID), nil)
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	return nil
}
