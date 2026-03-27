package adapters

import (
	"context"
	"testing"
	"time"

	"okapi/internal/models"
)

type mockAdapter struct {
	id string
}

func (m *mockAdapter) ID() string                  { return m.id }
func (m *mockAdapter) DisplayName() string         { return m.id }
func (m *mockAdapter) PollInterval() time.Duration { return time.Minute }
func (m *mockAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	return &models.StatusResponse{Service: m.id}, nil
}

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	a := &mockAdapter{id: "test"}

	r.Register(a)

	got, ok := r.Get("test")
	if !ok {
		t.Fatal("expected adapter not found")
	}
	if got.ID() != "test" {
		t.Errorf("expected test, got %s", got.ID())
	}

	all := r.All()
	if len(all) != 1 {
		t.Errorf("expected 1 adapter, got %d", len(all))
	}
}
