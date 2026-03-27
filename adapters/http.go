package adapters

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"okapi/internal/models"
)

type GenericHTTPAdapter struct {
	cfg        StatuspageConfig
	httpClient *http.Client
}

func NewGenericHTTPAdapter(cfg StatuspageConfig) *GenericHTTPAdapter {
	return &GenericHTTPAdapter{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *GenericHTTPAdapter) ID() string          { return a.cfg.ID }
func (a *GenericHTTPAdapter) DisplayName() string { return a.cfg.DisplayName }
func (a *GenericHTTPAdapter) PollInterval() time.Duration {
	return time.Duration(a.cfg.PollIntervalSeconds) * time.Second
}

func (a *GenericHTTPAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.cfg.Subdomain, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	status := models.StatusOperational
	summary := fmt.Sprintf("Service is reachable (Status: %d)", resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		status = models.StatusMajorOutage
		summary = fmt.Sprintf("Service returned non-success status: %d", resp.StatusCode)
	}

	return &models.StatusResponse{
		Service:    a.cfg.ID,
		Status:     status,
		Summary:    summary,
		FetchedAt:  time.Now(),
		DataSource: "http_check",
		SourceURL:  a.cfg.Subdomain,
	}, nil
}
