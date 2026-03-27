package api

import (
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"okapi/adapters"
	"okapi/api/middleware"
	"okapi/internal/cache"
	"okapi/internal/config"
	"okapi/internal/webhooks"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
)

func NewRouter(cfg *config.Config, registry *adapters.Registry, cache cache.Cache, w *webhooks.Manager, staticDist fs.FS) *chi.Mux {
	r := chi.NewRouter()
	handlers := NewHandlers(registry, cache, w, cfg)

	// Apply global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.MaxBodySize(1048576)) // 1MB limit
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Rate Limiter Middleware (e.g., 100 requests per minute per IP)
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   100,
		Interval: time.Minute,
	})
	if err != nil {
		log.Fatalf("failed to create rate limiter store: %v", err)
	}

	limiter, err := httplimit.NewMiddleware(store, httplimit.IPKeyFunc())
	if err != nil {
		log.Fatalf("failed to create rate limiter middleware: %v", err)
	}
	r.Use(limiter.Handle)

	// API Routes (Prefixed with /api)
	r.Route("/api", func(r chi.Router) {
		r.Get("/help", handlers.Help)
		r.Get("/_health", handlers.SelfHealth)

		// Routes that can be optionally protected
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(&cfg.Auth))

			r.Get("/maintenance", handlers.GetMaintenance)
			r.Get("/maintenance/{service}", handlers.GetMaintenance)

			r.Get("/incidents", handlers.GetRecentIncidents)
			r.Get("/incidents/{service}", handlers.GetRecentIncidents)

			r.Get("/health/{service}", handlers.GetServiceHealth)
			r.Get("/health", handlers.GetBatchHealth)
			r.Get("/services", handlers.ListServices)

			// Webhooks routes
			r.Get("/webhooks", handlers.ListWebhooks)
			r.Get("/webhooks/{service}", handlers.ListWebhooks)
			r.Post("/webhooks", handlers.RegisterWebhook)
			r.Delete("/webhooks/{id}", handlers.DeleteWebhook)
		})
	})

	// Serve static files if provided
	if staticDist != nil {
		fileServer := http.FileServer(http.FS(staticDist))

		// 1. Serve static files explicitly for folders
		r.Handle("/assets/*", fileServer)

		// 2. For individual files at root, we can use a helper or just let NotFound handle them
		r.Get("/favicon.ico", fileServer.ServeHTTP)
		r.Get("/manifest.json", fileServer.ServeHTTP)
		r.Get("/apple-touch-icon.png", fileServer.ServeHTTP)

		// 3. Fallback to index.html for SPA routes or any other unmatched path
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			// If it's a request for /api/* that actually 404ed, return JSON 404
			if strings.HasPrefix(r.URL.Path, "/api/") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": "API route not found"}`))
				return
			}

			// If it's a request for a file (has an extension), try to serve it from fileServer
			if strings.Contains(r.URL.Path, ".") {
				fileServer.ServeHTTP(w, r)
				return
			}

			// For any other path (like /services, /dashboard, /etc), serve index.html
			indexData, err := fs.ReadFile(staticDist, "index.html")
			if err != nil {
				// If index.html is missing, the build is broken
				http.Error(w, "Dashboard not found", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(indexData)
		})
	}

	return r
}
