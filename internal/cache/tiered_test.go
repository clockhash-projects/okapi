package cache

import (
	"context"
	"testing"
	"time"

	"okapi/internal/models"
)

func TestTieredCache(t *testing.T) {
	ctx := context.Background()
	l1 := NewMemoryCache()
	l2 := NewMemoryCache()
	localTTL := 1 * time.Second
	c := NewTieredCache(l1, l2, localTTL)

	tests := []struct {
		name    string
		action  func() error
		check   func() error
		wantErr bool
	}{
		{
			name: "Set updates both tiers",
			action: func() error {
				return c.Set(ctx, "k1", &models.StatusResponse{Service: "s1"}, 10*time.Second)
			},
			check: func() error {
				v1, _ := l1.Get(ctx, "k1")
				v2, _ := l2.Get(ctx, "k1")
				if v1 == nil || v2 == nil {
					t.Error("expected value in both tiers")
				}
				return nil
			},
		},
		{
			name: "Get fills L1 from L2",
			action: func() error {
				// Put only in L2
				_ = l2.Set(ctx, "k2", &models.StatusResponse{Service: "s2"}, 10*time.Second)
				_, err := c.Get(ctx, "k2")
				return err
			},
			check: func() error {
				v1, _ := l1.Get(ctx, "k2")
				if v1 == nil {
					t.Error("expected L1 to be backfilled")
				}
				return nil
			},
		},
		{
			name: "GetAll fills L1 from L2",
			action: func() error {
				_ = l2.Set(ctx, "k3", &models.StatusResponse{Service: "s3"}, 10*time.Second)
				_ = l2.Set(ctx, "k4", &models.StatusResponse{Service: "s4"}, 10*time.Second)
				_, err := c.GetAll(ctx, []string{"k3", "k4"})
				return err
			},
			check: func() error {
				v3, _ := l1.Get(ctx, "k3")
				v4, _ := l1.Get(ctx, "k4")
				if v3 == nil || v4 == nil {
					t.Error("expected L1 to be backfilled from GetAll")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.action(); (err != nil) != tt.wantErr {
				t.Errorf("%s: action() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if tt.check != nil {
				if err := tt.check(); err != nil {
					t.Errorf("%s: check() error = %v", tt.name, err)
				}
			}
		})
	}
}
