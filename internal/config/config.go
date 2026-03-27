package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Auth    AuthConfig    `yaml:"auth"`
	Cache   CacheConfig   `yaml:"cache"`
	Polling PollingConfig `yaml:"polling"`
	History HistoryConfig `yaml:"history"`
	Log     LogConfig     `yaml:"log"`
	CORS    CORSConfig    `yaml:"cors"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type ServerConfig struct {
	Port                int `yaml:"port"`
	ReadTimeoutSeconds  int `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int `yaml:"write_timeout_seconds"`
}

type AuthConfig struct {
	Enabled bool     `yaml:"enabled"`
	APIKeys []string `yaml:"api_keys"`
}

type CacheConfig struct {
	Backend           string `yaml:"backend"`
	RedisURL          string `yaml:"redis_url"`
	DefaultTTLSeconds int    `yaml:"default_ttl_seconds"`
}

type PollingConfig struct {
	DefaultIntervalSeconds int `yaml:"default_interval_seconds"`
	JitterSeconds          int `yaml:"jitter_seconds"`
}

type HistoryConfig struct {
	Enabled        bool   `yaml:"enabled"`
	RetentionHours int    `yaml:"retention_hours"`
	Storage        string `yaml:"storage"` // postgres
	PostgresURL    string `yaml:"postgres_url"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func (c *Config) Validate() error {
	if c.History.Enabled && c.History.PostgresURL == "" {
		return fmt.Errorf("history is enabled but OKAPI_HISTORY_POSTGRES_URL is missing")
	}
	if c.Cache.Backend == "redis" && c.Cache.RedisURL == "" {
		return fmt.Errorf("cache backend is redis but OKAPI_CACHE_REDIS_URL is missing")
	}
	return nil
}

func Load(path string) (*Config, error) {
	// Try the provided path first
	f, err := os.Open(path)
	if err != nil {
		// Fallback to default local path
		path = "config.yaml"
		f, err = os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
	}
	defer func() { _ = f.Close() }()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Override with environment variables
	overrideWithEnv(&cfg)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}

func overrideWithEnv(cfg *Config) {
	if v := os.Getenv("OKAPI_SERVER_PORT"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = i
		}
	}
	if v := os.Getenv("OKAPI_AUTH_ENABLED"); v != "" {
		cfg.Auth.Enabled = v == "true"
	}
	if v := os.Getenv("OKAPI_AUTH_API_KEYS"); v != "" {
		cfg.Auth.APIKeys = strings.Split(v, ",")
	}
	if v := os.Getenv("OKAPI_CACHE_BACKEND"); v != "" {
		cfg.Cache.Backend = v
	}
	if v := os.Getenv("OKAPI_CACHE_REDIS_URL"); v != "" {
		cfg.Cache.RedisURL = v
	}
	if v := os.Getenv("OKAPI_HISTORY_STORAGE"); v != "" {
		cfg.History.Storage = v
	}
	if v := os.Getenv("OKAPI_HISTORY_POSTGRES_URL"); v != "" {
		cfg.History.PostgresURL = v
	}
	if v := os.Getenv("OKAPI_CORS_ALLOWED_ORIGINS"); v != "" {
		cfg.CORS.AllowedOrigins = strings.Split(v, ",")
	}
}
