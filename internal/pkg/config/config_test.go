package config

import (
    "strings"
    "testing"
    "time"
)

func TestParseConfig(t *testing.T) {
    yamlData := `
provider:
  type: open-meteo
location:
  lat: 53.6688
  long: 23.8223
cache:
  enabled: true
  ttl: 5m
`

    reader := strings.NewReader(yamlData)
    cfg, err := Parse(reader)

    if err != nil {
        t.Errorf("Parse() returned error: %v", err)
    }

    if cfg.Provider.Type != "open-meteo" {
        t.Errorf("Expected provider type 'open-meteo', got '%s'", cfg.Provider.Type)
    }

    if cfg.Location.Lat != 53.6688 {
        t.Errorf("Expected lat 53.6688, got %f", cfg.Location.Lat)
    }

    if cfg.Location.Long != 23.8223 {
        t.Errorf("Expected long 23.8223, got %f", cfg.Location.Long)
    }

    if cfg.Cache.Enabled != true {
        t.Error("Expected cache enabled true")
    }

    if cfg.Cache.TTL != 5*time.Minute {
        t.Errorf("Expected TTL 5m, got %v", cfg.Cache.TTL)
    }
}

func TestParseConfigDefaultTTL(t *testing.T) {
    yamlData := `
provider:
  type: open-meteo
location:
  lat: 53.6688
  long: 23.8223
cache:
  enabled: true
`

    reader := strings.NewReader(yamlData)
    cfg, err := Parse(reader)

    if err != nil {
        t.Errorf("Parse() returned error: %v", err)
    }

    if cfg.Cache.TTL != 5*time.Minute {
        t.Errorf("Expected default TTL 5m, got %v", cfg.Cache.TTL)
    }
}

func TestParseConfigCacheDisabled(t *testing.T) {
    yamlData := `
provider:
  type: open-meteo
location:
  lat: 53.6688
  long: 23.8223
cache:
  enabled: false
  ttl: 5m
`

    reader := strings.NewReader(yamlData)
    cfg, err := Parse(reader)

    if err != nil {
        t.Errorf("Parse() returned error: %v", err)
    }

    if cfg.Cache.Enabled != false {
        t.Error("Expected cache enabled false")
    }
}

func TestParseConfigInvalidYAML(t *testing.T) {
    yamlData := `
provider:
  type: open-meteo
location:
  lat: invalid
  long: 23.8223
`

    reader := strings.NewReader(yamlData)
    _, err := Parse(reader)

    if err == nil {
        t.Error("Parse() should return error for invalid YAML")
    }
}