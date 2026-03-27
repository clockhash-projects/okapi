package code

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"okapi/internal/models"
)

type HerokuAdapter struct {
	BaseURL string
}

func (a *HerokuAdapter) ID() string          { return "heroku" }
func (a *HerokuAdapter) DisplayName() string { return "Heroku" }
func (a *HerokuAdapter) PollInterval() time.Duration {
	return 60 * time.Second
}

func (a *HerokuAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	url := "https://status.heroku.com/api/v4/current-status"
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
		Status []struct {
			System string `json:"system"`
			Status string `json:"status"`
		} `json:"status"`
		Incidents []struct {
			ID          int       `json:"id"`
			Title       string    `json:"title"`
			State       string    `json:"state"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
			FullMessage string    `json:"full_message"`
		} `json:"incidents"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	status := models.StatusOperational
	summary := "All systems operational"
	incidents := []models.Incident{}

	for _, s := range data.Status {
		if s.Status != "green" {
			status = models.StatusDegraded
			summary = "Issues detected on Heroku systems"
			break
		}
	}

	for _, i := range data.Incidents {
		incidents = append(incidents, models.Incident{
			ID:        fmt.Sprintf("%d", i.ID),
			Title:     i.Title,
			Status:    i.State,
			Body:      i.FullMessage,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
		})
	}

	return &models.StatusResponse{
		Service:    "heroku",
		Status:     status,
		Summary:    summary,
		Incidents:  incidents,
		FetchedAt:  time.Now(),
		DataSource: "official_api",
		SourceURL:  "https://status.heroku.com",
	}, nil
}
