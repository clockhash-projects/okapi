package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"syscall"

	"okapi/adapters"
	"okapi/internal/cache"
	"okapi/internal/config"
	"okapi/internal/logger"
	"okapi/internal/models"
	"okapi/internal/webhooks"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handlers struct {
	registry *adapters.Registry
	cache    cache.Cache
	webhooks *webhooks.Manager
	cfg      *config.Config
}

func NewHandlers(r *adapters.Registry, c cache.Cache, w *webhooks.Manager, cfg *config.Config) *Handlers {
	return &Handlers{
		registry: r,
		cache:    c,
		webhooks: w,
		cfg:      cfg,
	}
}

func (h *Handlers) RegisterWebhook(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		sendError(w, "feature_disabled", "Webhooks are not enabled", http.StatusNotImplemented)
		return
	}

	var req struct {
		URL      string   `json:"url"`
		Services []string `json:"services"`
		Secret   string   `json:"secret"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid_request", "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		sendError(w, "invalid_request", "URL is required", http.StatusUnprocessableEntity)
		return
	}

	id := h.webhooks.Register(webhooks.Webhook{
		URL:      req.URL,
		Services: req.Services,
		Secret:   req.Secret,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"id": id,
	}); err != nil {
		logger.Log.Error("failed to encode response", zap.Error(err))
	}
}

func (h *Handlers) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if h.webhooks == nil {
		sendError(w, "feature_disabled", "Webhooks are not enabled", http.StatusNotImplemented)
		return
	}

	h.webhooks.Delete(id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetRecentIncidents(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service")

	allServices := h.registry.All()
	if len(allServices) == 0 {
		sendError(w, "no_services", "No services registered", http.StatusNotFound)
		return
	}

	serviceIDs := make([]string, len(allServices))
	for i, adapter := range allServices {
		serviceIDs[i] = adapter.ID()
	}

	cachedResponses, err := h.cache.GetAll(r.Context(), serviceIDs)
	if err != nil {
		sendError(w, "internal_error", "Failed to fetch incidents from cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var incidents []models.Incident
	for _, res := range cachedResponses {
		if serviceID == "" || serviceID == res.Service {
			incidents = append(incidents, res.Incidents...)
		}
	}

	limit, offset := getPagination(r)
	if limit == 0 {
		limit = 50
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"incidents": incidents,
		"count":     len(incidents),
		"limit":     limit,
		"offset":    offset,
	}); err != nil {
		logger.Log.Error("failed to encode response", zap.Error(err))
	}
}

func (h *Handlers) Help(w http.ResponseWriter, r *http.Request) {
	help := map[string]interface{}{
		"name":        "Okapi API",
		"description": "Universal Service Health Proxy & API.",
		"endpoints": []map[string]string{
			{"method": "GET", "path": "/health/{service}", "description": "Current status for a single service (e.g., /health/github)"},
			{"method": "GET", "path": "/health", "description": "Batch status check (e.g., ?services=aws,github)"},
			{"method": "GET", "path": "/incidents", "description": "Consolidated view of all services currently experiencing issues"},
			{"method": "GET", "path": "/incidents/{service}", "description": "Recent incidents for a specific service"},
			{"method": "GET", "path": "/maintenance", "description": "Consolidated view of all upcoming scheduled maintenance"},
			{"method": "GET", "path": "/maintenance/{service}", "description": "Upcoming maintenance for a specific service"},
			{"method": "GET", "path": "/services", "description": "List all supported service IDs"},
			{"method": "GET", "path": "/_health", "description": "Self-health check for Okapi"},
			{"method": "GET", "path": "/webhooks", "description": "List all registered webhook subscriptions"},
			{"method": "POST", "path": "/webhooks", "description": "Register a new webhook for status changes"},
			{"method": "DELETE", "path": "/webhooks/{id}", "description": "Remove a registered webhook"},
			{"method": "GET", "path": "/help", "description": "Show this help information"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(help); err != nil {
		logger.Log.Error("failed to encode help", zap.Error(err))
	}
}

func (h *Handlers) GetMaintenance(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service")

	if serviceID != "" {
		adapter, ok := h.registry.Get(serviceID)
		if !ok {
			sendError(w, "service_not_found", "No adapter registered for service: "+serviceID, http.StatusNotFound)
			return
		}

		res, err := adapter.Fetch(r.Context())
		if err != nil {
			sendError(w, "upstream_error", "Failed to fetch status: "+err.Error(), http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"service":     res.Service,
			"maintenance": res.Maintenance,
		}); err != nil {
			logger.Log.Error("failed to encode maintenance", zap.Error(err))
		}
		return
	}

	allServices := h.registry.All()
	if len(allServices) == 0 {
		sendError(w, "no_services", "No services registered", http.StatusNotFound)
		return
	}

	serviceIDs := make([]string, len(allServices))
	for i, adapter := range allServices {
		serviceIDs[i] = adapter.ID()
	}

	cachedResponses, err := h.cache.GetAll(r.Context(), serviceIDs)
	if err != nil {
		sendError(w, "internal_error", "Failed to fetch maintenance from cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var maintenanceEvents []models.Maintenance
	for _, res := range cachedResponses {
		maintenanceEvents = append(maintenanceEvents, res.Maintenance...)
	}

	limit, offset := getPagination(r)
	if limit == 0 {
		limit = 100
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"scheduled_maintenance": maintenanceEvents,
		"count":                 len(maintenanceEvents),
		"limit":                 limit,
		"offset":                offset,
	}); err != nil {
		logger.Log.Error("failed to encode maintenance", zap.Error(err))
	}
}

func (h *Handlers) SelfHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"version": "1.0.0",
	}); err != nil {
		logger.Log.Error("failed to encode self-health", zap.Error(err))
	}
}

func (h *Handlers) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	if h.webhooks == nil {
		sendError(w, "feature_disabled", "Webhooks are not enabled", http.StatusNotImplemented)
		return
	}

	webhooks := h.webhooks.List()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(webhooks); err != nil {
		logger.Log.Error("failed to encode webhooks", zap.Error(err))
	}
}

func (h *Handlers) GetServiceHealth(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service")
	adapter, ok := h.registry.Get(serviceID)
	if !ok {
		sendError(w, "service_not_found", "No adapter registered for service: "+serviceID, http.StatusNotFound)
		return
	}

	var res *models.StatusResponse
	var err error

	if h.cache != nil {
		res, err = h.cache.Get(r.Context(), serviceID)
		if err == nil && res != nil {
			res.Cached = true
			w.Header().Set("Content-Type", "application/json")
			if encErr := json.NewEncoder(w).Encode(res); encErr != nil {
				logger.Log.Error("failed to encode health", zap.Error(encErr))
			}
			return
		}
	}

	res, err = adapter.Fetch(r.Context())
	if err != nil {
		sendError(w, "upstream_error", "Failed to fetch status: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Log.Error("failed to encode health", zap.Error(err))
	}
}

func (h *Handlers) GetBatchHealth(w http.ResponseWriter, r *http.Request) {
	servicesParam := r.URL.Query().Get("services")
	if servicesParam == "" {
		sendError(w, "invalid_request", "Missing 'services' query parameter", http.StatusBadRequest)
		return
	}

	ids := strings.Split(servicesParam, ",")
	results := make(map[string]*models.StatusResponse)

	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}

		adapter, ok := h.registry.Get(id)
		if !ok {
			continue
		}

		var res *models.StatusResponse
		var err error
		if h.cache != nil {
			res, err = h.cache.Get(r.Context(), id)
		}

		if err != nil || res == nil {
			res, err = adapter.Fetch(r.Context())
		}

		if err == nil && res != nil {
			results[id] = res
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
		"count":   len(results),
	}); err != nil {
		logEncodingError(err)
	}
}

func (h *Handlers) ListServices(w http.ResponseWriter, r *http.Request) {
	adapters := h.registry.All()
	services := make([]string, 0, len(adapters))
	for _, a := range adapters {
		services = append(services, a.ID())
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(services); err != nil {
		logEncodingError(err)
	}
}

func sendError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	}); err != nil {
		logger.Log.Error("failed to encode error", zap.Error(err))
	}
}

func logEncodingError(err error) {
	if err == nil {
		return
	}
	if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) {
		return
	}
	msg := err.Error()
	if strings.Contains(msg, "broken pipe") || strings.Contains(msg, "connection reset") {
		return
	}

	logger.Log.Error("failed to encode response", zap.Error(err))
}

func getPagination(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}
	return
}
