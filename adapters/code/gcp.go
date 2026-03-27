package code

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"okapi/internal/models"
)

type GCPAdapter struct {
	BaseURL string
}

func (a *GCPAdapter) ID() string          { return "gcp" }
func (a *GCPAdapter) DisplayName() string { return "Google Cloud Platform" }
func (a *GCPAdapter) PollInterval() time.Duration {
	return 60 * time.Second
}

func (a *GCPAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	url := "https://status.cloud.google.com/incidents.json"
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

	var data []struct {
		ID             string    `json:"id"`
		ServiceKey     string    `json:"service_key"`
		ServiceName    string    `json:"service_name"`
		Severity       string    `json:"severity"`
		StatusDescription string `json:"status_description"`
		Begin          time.Time `json:"begin"`
		End            string    `json:"end"`
		ExternalDesc   string    `json:"external_desc"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	status := models.StatusOperational
	summary := "All systems operational"
	incidents := []models.Incident{}

	activeCount := 0
	for _, item := range data {
		var endTime time.Time
		if item.End != "" {
			endTime, _ = time.Parse(time.RFC3339, item.End)
		}

		if endTime.IsZero() || endTime.After(time.Now()) {
			activeCount++
			incidents = append(incidents, models.Incident{
				ID:        item.ID,
				Title:     fmt.Sprintf("%s: %s", item.ServiceName, item.StatusDescription),
				Status:    item.Severity,
				Body:      item.ExternalDesc,
				CreatedAt: item.Begin,
			})
		}
	}

	if activeCount > 0 {
		status = models.StatusDegraded
		summary = fmt.Sprintf("%d active incidents reported", activeCount)
	}

	return &models.StatusResponse{
		Service:    "gcp",
		Status:     status,
		Summary:    summary,
		Incidents:  incidents,
		FetchedAt:  time.Now(),
		DataSource: "official_api",
		SourceURL:  "https://status.cloud.google.com",
	}, nil
}
