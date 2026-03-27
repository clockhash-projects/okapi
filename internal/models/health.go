package models

import (
	"encoding/json"
	"time"
)

// Status represents the normalized status of a service.
type Status string

const (
	StatusOperational   Status = "operational"
	StatusDegraded      Status = "degraded"
	StatusPartialOutage Status = "partial_outage"
	StatusMajorOutage   Status = "major_outage"
	StatusMaintenance   Status = "maintenance"
	StatusUnknown       Status = "unknown"
)

// StatusPoint represents a historical status at a specific time.
type StatusPoint struct {
	Status Status    `json:"status"`
	Time   time.Time `json:"time"`
}

// StatusResponse represents the unified response shape for all health endpoints.
// It contains core fields for system logic and a Metadata map for flexible, service-specific data.
type StatusResponse struct {
	Service    string         `json:"service"`
	Status     Status         `json:"status"`
	Summary    string         `json:"summary"`
	FetchedAt  time.Time      `json:"fetched_at"`
	DataSource string         `json:"data_source"`
	SourceURL  string         `json:"source_url,omitempty"`
	Cached     bool           `json:"cached,omitempty"`
	Components []Component    `json:"components,omitempty"`
	Incidents  []Incident     `json:"incidents,omitempty"`
	Maintenance []Maintenance `json:"scheduled_maintenance,omitempty"`
	History    []StatusPoint  `json:"recent_history,omitempty"`
	Metadata   map[string]any `json:"-"` // Metadata is flattened into the top-level JSON
}

// MarshalJSON implements json.Marshaler to flatten Metadata into the top-level object.
func (r StatusResponse) MarshalJSON() ([]byte, error) {
	type Alias StatusResponse
	b, err := json.Marshal(Alias(r))
	if err != nil {
		return nil, err
	}

	if len(r.Metadata) == 0 {
		return b, nil
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	for k, v := range r.Metadata {
		m[k] = v
	}

	return json.Marshal(m)
}

// Maintenance represents a scheduled maintenance event.
type Maintenance struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Summary   string    `json:"summary"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Component represents an individual part of a service.
type Component struct {
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Incident represents an ongoing or recent issue.
type Incident struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
