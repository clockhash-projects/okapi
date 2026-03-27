package adapters

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Registry struct {
	mu       sync.RWMutex
	adapters map[string]HealthAdapter
}

func NewRegistry() *Registry {
	return &Registry{
		adapters: make(map[string]HealthAdapter),
	}
}

func (r *Registry) Register(a HealthAdapter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[a.ID()] = a
}

func (r *Registry) Get(id string) (HealthAdapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.adapters[id]
	return a, ok
}

func (r *Registry) All() []HealthAdapter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var all []HealthAdapter
	for _, a := range r.adapters {
		all = append(all, a)
	}
	return all
}

func (r *Registry) LoadFromConfig(dir string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) && path == dir {
				return nil
			}
			return err
		}

		if !d.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read config %s: %w", path, err)
			}

			var base struct {
				ID   string `yaml:"id"`
				Kind string `yaml:"kind"`
			}
			if err := yaml.Unmarshal(data, &base); err != nil {
				return fmt.Errorf("failed to unmarshal config base %s: %w", path, err)
			}

			// Validation
			if base.ID == "" {
				return fmt.Errorf("config %s: id is required", path)
			}
			if base.Kind == "" {
				return fmt.Errorf("config %s: kind is required", path)
			}

			switch base.Kind {
			case "statuspage":
				var cfg StatuspageConfig
				if err := yaml.Unmarshal(data, &cfg); err != nil {
					return err
				}
				if cfg.Subdomain == "" {
					return fmt.Errorf("config %s: subdomain is required for statuspage", path)
				}
				if cfg.DisplayName == "" {
					cfg.DisplayName = cfg.ID
				}
				r.Register(NewStatuspageAdapter(cfg))
			case "http":
				var cfg StatuspageConfig
				if err := yaml.Unmarshal(data, &cfg); err != nil {
					return err
				}
				if cfg.Subdomain == "" {
					return fmt.Errorf("config %s: url (subdomain field) is required for http", path)
				}
				if cfg.DisplayName == "" {
					cfg.DisplayName = cfg.ID
				}
				r.Register(NewGenericHTTPAdapter(cfg))
			case "rss":
				var cfg RSSConfig
				if err := yaml.Unmarshal(data, &cfg); err != nil {
					return err
				}
				if cfg.URL == "" {
					return fmt.Errorf("config %s: url is required for rss", path)
				}
				if cfg.DisplayName == "" {
					cfg.DisplayName = cfg.ID
				}
				r.Register(NewRSSAdapter(cfg))
			case "statusio":
				var cfg StatuspageConfig
				if err := yaml.Unmarshal(data, &cfg); err != nil {
					return err
				}
				if cfg.StatusioID == "" {
					return fmt.Errorf("config %s: statusio_id is required for statusio", path)
				}
				if cfg.DisplayName == "" {
					cfg.DisplayName = cfg.ID
				}
				r.Register(NewStatusioAdapter(cfg))
			default:
				return fmt.Errorf("config %s: unknown kind %s", path, base.Kind)
			}
		}
		return nil
	})
}
