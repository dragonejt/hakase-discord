package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
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

func ReadAssignment(span *sentry.Span, assignmentID string) (Assignment, error) {
	span = span.StartChild("readAssignment")
	defer span.Finish()

	assignment := Assignment{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/assignments?id=%s", settings.BACKEND_URL, assignmentID), nil)
	if err != nil {
		return assignment, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

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

func ListAssignments(span *sentry.Span, courseID string) ([]Assignment, error) {
	span = span.StartChild("listAssignments")
	defer span.Finish()

	assignments := []Assignment{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/assignments?course_id=%s", settings.BACKEND_URL, courseID), nil)
	if err != nil {
		return assignments, fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	client := span.GetTransaction().Context().Value(DiscordSession{}).(*discordgo.Session).Client
	response, err := client.Do(request)
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

func CreateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {
	span = span.StartChild("createAssignment")
	defer span.Finish()

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
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	client := span.GetTransaction().Context().Value(DiscordSession{}).(*discordgo.Session).Client
	response, err := client.Do(request)
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

func UpdateAssignment(span *sentry.Span, assignment Assignment) (Assignment, error) {
	span = span.StartChild("updateAssignment")
	defer span.Finish()

	jsonBody, err := json.Marshal(assignment)
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
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	client := span.GetTransaction().Context().Value(DiscordSession{}).(*discordgo.Session).Client
	response, err := client.Do(request)
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

func DeleteAssignment(span *sentry.Span, assignmentID string) error {
	span = span.StartChild("deleteAssignment")
	defer span.Finish()

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/assignments?id=%s", settings.BACKEND_URL, assignmentID), nil)
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.BACKEND_API_KEY))
	request.Header.Add(sentry.SentryTraceHeader, sentry.CurrentHub().GetTraceparent())
	request.Header.Add(sentry.SentryBaggageHeader, sentry.CurrentHub().GetBaggage())

	client := span.GetTransaction().Context().Value(DiscordSession{}).(*discordgo.Session).Client
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	return nil
}
