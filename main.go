package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"okapi/adapters"
	"okapi/adapters/code"
	"okapi/api"
	"okapi/internal/cache"
	"okapi/internal/config"
	"okapi/internal/logger"
	"okapi/internal/polling"
	"okapi/internal/webhooks"

	"go.uber.org/zap"
)

//go:embed dashboard/dist
var dashboardFS embed.FS

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		fmt.Println("Migration completed (no-op as history is removed).")
		os.Exit(0)
	}

	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	
	// Ensure logs are flushed before exit
	defer func() {
		_ = logger.Log.Sync()
	}()

	// Helper to handle fatal startup errors
	fatal := func(msg string, fields ...zap.Field) {
		logger.Log.Error(msg, fields...)
		// Do not Sync() here, the defer will handle it
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// cancel is called explicitly during shutdown; no defer here to avoid
	// a redundant cancel racing with the shutdown goroutine.

	// Initialize cache
	var c cache.Cache
	switch cfg.Cache.Backend {
	case "redis":
		c, err = cache.NewRedisCache(cfg.Cache.RedisURL)
		if err != nil {
			fatal("failed to initialize redis cache", zap.Error(err))
		}
	case "memory":
		c = cache.NewMemoryCache()
	default:
		fatal("no valid cache backend configured", zap.String("backend", cfg.Cache.Backend))
	}

	// Initialize webhooks
	w := webhooks.NewManager()

	// Initialize registry
	registry := adapters.NewRegistry()

	registry.Register(&code.AWSAdapter{})
	registry.Register(&code.AzureAdapter{})
	registry.Register(&code.GCPAdapter{})
	registry.Register(&code.HerokuAdapter{})
	registry.Register(&code.SlackAdapter{})

	if err = registry.LoadFromConfig("adapters/config"); err != nil {
		logger.Log.Fatal("failed to load YAML adapters", zap.Error(err))
	}

	worker := polling.NewWorker(registry, c, w, cfg)
	worker.Start(ctx)

	// Prepare dashboard static files
	subFS, err := fs.Sub(dashboardFS, "dashboard/dist")
	if err != nil {
		fatal("failed to create dashboard sub-filesystem", zap.Error(err))
	}
	r := api.NewRouter(cfg, registry, c, w, subFS)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSeconds) * time.Second,
	}

	go handleShutdown(srv, worker, cancel)

	logger.Log.Info("Starting Okapi", zap.String("addr", addr))
	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("failed to start server", zap.Error(err))
	}
}

// handleShutdown waits for an OS signal, gracefully drains the HTTP server,
// then cancels the worker context and waits for all workers to finish.
func handleShutdown(srv *http.Server, worker *polling.Worker, cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Log.Info("Shutting down Okapi...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if shutdownErr := srv.Shutdown(shutdownCtx); shutdownErr != nil {
		logger.Log.Error("HTTP server shutdown error", zap.Error(shutdownErr))
	}

	// Stop workers only after the HTTP server has stopped accepting requests.
	cancel()

	logger.Log.Info("Waiting for workers to finish...")
	worker.Wait()
	logger.Log.Info("Shutdown complete.")
}
