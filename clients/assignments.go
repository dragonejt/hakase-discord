package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/dragonejt/hakase-discord/settings"
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

func ReadAssignment(assignmentID string) (Assignment, error) {
	sentry.StartSpan(context.TODO(), "readAssignment")
	assignment := Assignment{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/assignments?id=%s", settings.BACKEND_URL, assignmentID), nil)
	if err != nil {
		return assignment, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))

	response, err := http.DefaultClient.Do(request)
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

func ListAssignments(courseID string) ([]Assignment, error) {
	sentry.StartSpan(context.TODO(), "listAssignments")
	assignments := []Assignment{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/assignments?course_id=%s", settings.BACKEND_URL, courseID), nil)
	if err != nil {
		return assignments, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))

	response, err := http.DefaultClient.Do(request)
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

func CreateAssignment(assignment Assignment) (Assignment, error) {
	sentry.StartSpan(context.TODO(), "createAssignment")
	jsonBody, err := json.Marshal(assignment)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to marshal assignment: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/assignments", settings.BACKEND_URL), bytes.NewReader(jsonBody))
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))

	response, err := http.DefaultClient.Do(request)
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

func UpdateAssignment(assignment Assignment) (Assignment, error) {
	sentry.StartSpan(context.TODO(), "updateAssignment")
	jsonBody, err := json.Marshal(assignment)
	slog.Info(string(jsonBody))
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to marshal assignment: %w", err)
	}

	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/assignments", settings.BACKEND_URL), bytes.NewReader(jsonBody))
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))

	response, err := http.DefaultClient.Do(request)
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

func DeleteAssignment(assignmentID string) error {
	sentry.StartSpan(context.TODO(), "deleteAssignment")
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/assignments?id=%s", settings.BACKEND_URL, assignmentID), nil)
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
