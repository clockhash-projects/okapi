package webhooks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"okapi/internal/models"
)

type Webhook struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Services  []string  `json:"services"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Payload struct {
	Event          string    `json:"event"`
	Service        string    `json:"service"`
	PreviousStatus string    `json:"previous_status"`
	CurrentStatus  string    `json:"current_status"`
	Summary        string    `json:"summary"`
	OccurredAt     time.Time `json:"occurred_at"`
}

type Manager struct {
	mu       sync.RWMutex
	webhooks map[string]Webhook
}

func NewManager() *Manager {
	return &Manager{
		webhooks: make(map[string]Webhook),
	}
}

func (m *Manager) Register(w Webhook) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if w.ID == "" {
		w.ID = fmt.Sprintf("hook_%d", time.Now().UnixNano())
	}
	w.CreatedAt = time.Now()
	m.webhooks[w.ID] = w
	return w.ID
}

func (m *Manager) Delete(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.webhooks, id)
}

func (m *Manager) List() []Webhook {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hooks := make([]Webhook, 0, len(m.webhooks))
	for _, w := range m.webhooks {
		hooks = append(hooks, w)
	}
	return hooks
}

func (m *Manager) GetByService(service string) []Webhook {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hooks := make([]Webhook, 0)
	for _, w := range m.webhooks {
		if m.shouldNotify(w, service) {
			hooks = append(hooks, w)
		}
	}
	return hooks
}

func (m *Manager) Notify(oldStatus, newStatus *models.StatusResponse) {
	if oldStatus.Status == newStatus.Status {
		return
	}

	payload := Payload{
		Event:          "status_changed",
		Service:        newStatus.Service,
		PreviousStatus: string(oldStatus.Status),
		CurrentStatus:  string(newStatus.Status),
		Summary:        newStatus.Summary,
		OccurredAt:     newStatus.FetchedAt,
	}

	m.mu.RLock()
	hooks := make([]Webhook, 0, len(m.webhooks))
	for _, w := range m.webhooks {
		hooks = append(hooks, w)
	}
	m.mu.RUnlock()

	for _, hook := range hooks {
		if m.shouldNotify(hook, newStatus.Service) {
			go func(h Webhook) {
				if err := m.send(h, payload); err != nil {
					log.Printf("failed to send webhook: %v", err)
				}
			}(hook)
		}
	}
}

func (m *Manager) shouldNotify(w Webhook, service string) bool {
	if service == "" {
		return true // Requesting all webhooks
	}
	if len(w.Services) == 0 {
		return true // Webhook watching all services
	}
	for _, s := range w.Services {
		if s == service {
			return true
		}
	}
	return false
}

func (m *Manager) send(w Webhook, p Payload) error {
	body, _ := json.Marshal(p)
	req, err := http.NewRequest("POST", w.URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if w.Secret != "" {
		mac := hmac.New(sha256.New, []byte(w.Secret))
		mac.Write(body)
		signature := hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Okapi-Signature", signature)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	return nil
}
