package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"okapi/internal/models"
)

type StatusioAdapter struct {
	cfg        StatuspageConfig
	httpClient *http.Client
}

func NewStatusioAdapter(cfg StatuspageConfig) *StatusioAdapter {
	return &StatusioAdapter{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *StatusioAdapter) ID() string          { return a.cfg.ID }
func (a *StatusioAdapter) DisplayName() string { return a.cfg.DisplayName }
func (a *StatusioAdapter) PollInterval() time.Duration {
	return time.Duration(a.cfg.PollIntervalSeconds) * time.Second
}

func (a *StatusioAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	// Status.io public status API
	url := fmt.Sprintf("https://api.status.io/v1/status/%s", a.cfg.StatusioID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream returned status: %d", resp.StatusCode)
	}

	var data struct {
		Result struct {
			StatusOverall struct {
				Status      string `json:"status"`
				UpdatedTime string `json:"updated_time"`
			} `json:"status_overall"`
			Incidents []struct {
				ID          string `json:"_id"`
				Title       string `json:"title"`
				Status      string `json:"status"`
				CreatedTime string `json:"created_time"`
			} `json:"incidents"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	status := models.StatusOperational
	switch data.Result.StatusOverall.Status {
	case "Operational":
		status = models.StatusOperational
	case "Degraded Performance":
		status = models.StatusDegraded
	case "Partial Service Disruption":
		status = models.StatusPartialOutage
	case "Service Disruption":
		status = models.StatusMajorOutage
	}

	pageID := a.cfg.StatusioPageID
	if pageID == "" {
		pageID = a.cfg.StatusioID
	}

	res := &models.StatusResponse{
		Service:    a.cfg.ID,
		Status:     status,
		Summary:    data.Result.StatusOverall.Status,
		FetchedAt:  time.Now(),
		DataSource: "status_io",
		SourceURL:  fmt.Sprintf("https://status.io/pages/%s", pageID),
	}

	for _, i := range data.Result.Incidents {
		createdAt, _ := time.Parse(time.RFC3339, i.CreatedTime)
		res.Incidents = append(res.Incidents, models.Incident{
			ID:        i.ID,
			Title:     i.Title,
			Status:    i.Status,
			CreatedAt: createdAt,
		})
	}

	return res, nil
}
