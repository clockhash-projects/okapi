package polling

import (
	"context"
	"sync"
	"time"

	"okapi/adapters"
	"okapi/internal/cache"
	"okapi/internal/config"
	"okapi/internal/logger"
	"okapi/internal/models"
	"okapi/internal/webhooks"

	"go.uber.org/zap"
)

type Worker struct {
	registry *adapters.Registry
	cache    cache.Cache
	webhooks *webhooks.Manager
	cfg      *config.Config
	wg       sync.WaitGroup
}

func NewWorker(r *adapters.Registry, c cache.Cache, w *webhooks.Manager, cfg *config.Config) *Worker {
	return &Worker{
		registry: r,
		cache:    c,
		webhooks: w,
		cfg:      cfg,
	}
}

func (w *Worker) Start(ctx context.Context) {
	adapters := w.registry.All()
	for _, a := range adapters {
		w.wg.Add(1)
		go w.pollAdapter(ctx, a)
	}
}

func (w *Worker) Wait() {
	w.wg.Wait()
}

func (w *Worker) pollAdapter(ctx context.Context, a adapters.HealthAdapter) {
	defer w.wg.Done()
	interval := a.PollInterval()
	if interval == 0 {
		interval = time.Duration(w.cfg.Polling.DefaultIntervalSeconds) * time.Second
	}

	// Initial poll
	w.doPoll(ctx, a)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.doPoll(ctx, a)
		}
	}
}

func (w *Worker) doPoll(ctx context.Context, a adapters.HealthAdapter) {
	start := time.Now()

	// Apply a per-call timeout as a safety measure
	fetchCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Fetch current status
	res, err := a.Fetch(fetchCtx)
	duration := time.Since(start).Seconds()
	logger.Log.Info("polled service", zap.String("service", a.ID()), zap.Float64("duration_s", duration))

	if err != nil {
		logger.Log.Error("polling error", zap.String("service", a.ID()), zap.Error(err))
		logger.PollDuration.WithLabelValues(a.ID(), "false").Observe(duration)

		// Create an error response to propagate the failure state
		res = &models.StatusResponse{
			Service:   a.ID(),
			Status:    models.StatusUnknown,
			Summary:   "Polling failed: " + err.Error(),
			FetchedAt: time.Now(),
		}
	} else {
		logger.PollDuration.WithLabelValues(a.ID(), "true").Observe(duration)
		logger.RecordStatus(a.ID(), res.Status)
	}

	// Check for status change if webhooks are enabled
	if w.webhooks != nil && w.cache != nil {
		// Try to get previous status from cache
		prev, err := w.cache.Get(ctx, a.ID())
		if err == nil && prev != nil {
			if prev.Status != res.Status {
				logger.Log.Info("status change detected",
					zap.String("service", a.ID()),
					zap.String("prev", string(prev.Status)),
					zap.String("curr", string(res.Status)),
				)
				w.webhooks.Notify(prev, res)
			}
		}
	}

	if w.cache != nil {
		ttl := time.Duration(w.cfg.Cache.DefaultTTLSeconds) * time.Second
		_ = w.cache.Set(ctx, a.ID(), res, ttl)
	}
}
