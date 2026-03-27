package webhooks

import (
	"testing"
)

func TestManager_Register_Delete(t *testing.T) {
	m := NewManager()

	w := Webhook{
		URL:      "http://example.com",
		Services: []string{"test"},
		Secret:   "secret",
	}

	id := m.Register(w)
	if id == "" {
		t.Fatal("expected non-empty ID")
	}

	hooks := m.List()
	if len(hooks) != 1 {
		t.Errorf("expected 1 hook, got %d", len(hooks))
	}

	m.Delete(id)
	hooks = m.List()
	if len(hooks) != 0 {
		t.Errorf("expected 0 hooks, got %d", len(hooks))
	}
}

func TestManager_GetByService(t *testing.T) {
	m := NewManager()

	m.Register(Webhook{URL: "h1", Services: []string{"s1"}})
	m.Register(Webhook{URL: "h2", Services: []string{"s2"}})
	m.Register(Webhook{URL: "h3", Services: []string{"s1", "s2"}})
	m.Register(Webhook{URL: "h4", Services: []string{}}) // All services

	s1Hooks := m.GetByService("s1")
	if len(s1Hooks) != 3 { // h1, h3, h4
		t.Errorf("expected 3 hooks for s1, got %d", len(s1Hooks))
	}

	allHooks := m.GetByService("")
	if len(allHooks) != 4 {
		t.Errorf("expected 4 hooks for all, got %d", len(allHooks))
	}
}
