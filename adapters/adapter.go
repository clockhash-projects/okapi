package adapters

import (
	"context"
	"time"

	"okapi/internal/models"
)

// HealthAdapter is the interface that all status adapters must implement.
type HealthAdapter interface {
	ID() string
	DisplayName() string
	PollInterval() time.Duration
	Fetch(ctx context.Context) (*models.StatusResponse, error)
}
