package code

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"okapi/internal/models"
)

type SlackAdapter struct {
	BaseURL string
}

func (a *SlackAdapter) ID() string          { return "slack" }
func (a *SlackAdapter) DisplayName() string { return "Slack" }
func (a *SlackAdapter) PollInterval() time.Duration {
	return 60 * time.Second
}

func (a *SlackAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	url := "https://slack-status.com/api/v2.0.0/current"
	if a.BaseURL != "" {
		url = a.BaseURL
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data struct {
		Status    string `json:"status"`
		DateCreated time.Time `json:"date_created"`
		ActiveIncidents []struct {
			ID          int       `json:"id"`
			Title       string    `json:"title"`
			Type        string    `json:"type"`
			Status      string    `json:"status"`
			DateCreated time.Time `json:"date_created"`
			DateUpdated time.Time `json:"date_updated"`
		} `json:"active_incidents"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	status := models.StatusOperational
	summary := "Slack is up and running"

	switch data.Status {
	case "ok":
		status = models.StatusOperational
	case "degraded":
		status = models.StatusDegraded
	case "outage":
		status = models.StatusMajorOutage
	}

	incidents := []models.Incident{}
	for _, i := range data.ActiveIncidents {
		incidents = append(incidents, models.Incident{
			ID:        fmt.Sprintf("%d", i.ID),
			Title:     i.Title,
			Status:    i.Status,
			CreatedAt: i.DateCreated,
			UpdatedAt: i.DateUpdated,
		})
	}

	return &models.StatusResponse{
		Service:    "slack",
		Status:     status,
		Summary:    summary,
		Incidents:  incidents,
		FetchedAt:  time.Now(),
		DataSource: "official_api",
		SourceURL:  "https://status.slack.com",
	}, nil
}
