package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"okapi/internal/logger"
	"okapi/internal/models"

	"go.uber.org/zap"
)

type StatuspageConfig struct {
	ID                  string            `yaml:"id"`
	DisplayName         string            `yaml:"display_name"`
	Kind                string            `yaml:"kind"`
	Subdomain           string            `yaml:"subdomain"`
	StatusioID          string            `yaml:"statusio_id"`
	StatusioPageID      string            `yaml:"statusio_page_id"`
	PollIntervalSeconds int               `yaml:"poll_interval_seconds"`
	ComponentAliases    map[string]string `yaml:"component_aliases"`
}

type StatuspageAdapter struct {
	cfg        StatuspageConfig
	httpClient *http.Client
}

func NewStatuspageAdapter(cfg StatuspageConfig) *StatuspageAdapter {
	return &StatuspageAdapter{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *StatuspageAdapter) ID() string          { return a.cfg.ID }
func (a *StatuspageAdapter) DisplayName() string { return a.cfg.DisplayName }
func (a *StatuspageAdapter) PollInterval() time.Duration {
	return time.Duration(a.cfg.PollIntervalSeconds) * time.Second
}

func (a *StatuspageAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	scheme := "https"
	if strings.Contains(a.cfg.Subdomain, "127.0.0.1") || strings.Contains(a.cfg.Subdomain, "localhost") {
		scheme = "http"
	}
	url := fmt.Sprintf("%s://%s/api/v2/summary.json", scheme, a.cfg.Subdomain)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream returned status: %d", resp.StatusCode)
	}

	var data struct {
		Page struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"page"`
		Status struct {
			Indicator   string `json:"indicator"`
			Description string `json:"description"`
		} `json:"status"`
		Components []struct {
			Name      string    `json:"name"`
			Status    string    `json:"status"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"components"`
		Incidents []struct {
			ID        string    `json:"id"`
			Name      string    `json:"name"`
			Status    string    `json:"status"`
			Body      string    `json:"shortlink"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"incidents"`
		Maintenance []struct {
			ID             string    `json:"id"`
			Name           string    `json:"name"`
			Status         string    `json:"status"`
			UpdatedAt      time.Time `json:"updated_at"`
			ScheduledFor   time.Time `json:"scheduled_for"`
			ScheduledUntil time.Time `json:"scheduled_until"`
			Updates        []struct {
				Body string `json:"body"`
			} `json:"incident_updates"`
		} `json:"scheduled_maintenances"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	res := &models.StatusResponse{
		Service:    a.cfg.ID,
		Status:     mapStatuspageStatus(data.Status.Indicator),
		Summary:    cleanString(data.Status.Description),
		FetchedAt:  time.Now(),
		DataSource: "official_api",
		SourceURL:  url,
	}

	for _, c := range data.Components {
		res.Components = append(res.Components, models.Component{
			Name:      c.Name,
			Status:    mapStatuspageStatus(c.Status),
			UpdatedAt: c.UpdatedAt,
		})
	}

	for _, i := range data.Incidents {
		res.Incidents = append(res.Incidents, models.Incident{
			ID:        i.ID,
			Title:     i.Name,
			Status:    i.Status,
			Body:      i.Body,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
		})
	}

	for _, m := range data.Maintenance {
		summary := ""
		if len(m.Updates) > 0 {
			summary = cleanString(m.Updates[0].Body)
		}
		res.Maintenance = append(res.Maintenance, models.Maintenance{
			ID:        m.ID,
			Title:     m.Name,
			Status:    m.Status,
			Summary:   summary,
			StartsAt:  m.ScheduledFor,
			EndsAt:    m.ScheduledUntil,
			UpdatedAt: m.UpdatedAt,
		})
	}

	return res, nil
}

func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	return strings.TrimSpace(s)
}

func mapStatuspageStatus(s string) models.Status {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "", "none", "operational", "resolved":
		return models.StatusOperational
	case "minor", "degraded_performance", "investigating":
		return models.StatusDegraded
	case "partial", "partial_outage", "identified":
		return models.StatusPartialOutage
	case "major", "major_outage", "monitoring":
		return models.StatusMajorOutage
	case "maintenance", "under_maintenance", "scheduled":
		return models.StatusMaintenance
	default:
		// Log unknown status for visibility if logger is initialized
		if s != "" && logger.Log != nil {
			logger.Log.Warn("Unknown statuspage indicator", zap.String("indicator", s))
		}
		return models.StatusUnknown
	}
}
