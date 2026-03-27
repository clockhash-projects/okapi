package config

import (
	"os"
	"testing"
)

func TestOverrideWithEnv(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 8080},
	}

	_ = os.Setenv("OKAPI_SERVER_PORT", "9090")
	defer func() { _ = os.Unsetenv("OKAPI_SERVER_PORT") }()

	_ = os.Setenv("OKAPI_CACHE_BACKEND", "redis")
	defer func() { _ = os.Unsetenv("OKAPI_CACHE_BACKEND") }()

	overrideWithEnv(cfg)

	if cfg.Server.Port != 9090 {
		t.Errorf("expected 9090, got %d", cfg.Server.Port)
	}
	if cfg.Cache.Backend != "redis" {
		t.Errorf("expected redis, got %s", cfg.Cache.Backend)
	}
}

func TestLoad(t *testing.T) {
	content := `
server:
  port: 7070
cache:
  backend: memory
`
	tmpfile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	_, err = tmpfile.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	err = tmpfile.Close()
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Port != 7070 {
		t.Errorf("expected 7070, got %d", cfg.Server.Port)
	}
}
