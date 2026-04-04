package main

import (
    "testing"
    "time"

    pogodaby "github.com/Kurome00/weather-app.git/internal/adapters/pogoda_by"
    "github.com/Kurome00/weather-app.git/internal/adapters/weather"
    "github.com/Kurome00/weather-app.git/internal/pkg/app/cli"
    "github.com/Kurome00/weather-app.git/internal/pkg/config"
)

type testLogger struct{}

func (testLogger) Info(string)  {}
func (testLogger) Debug(string) {}
func (testLogger) Error(string) {}

func TestGetProviderOpenMeteo(t *testing.T) {
    cfg := config.Config{
        Provider: config.Provider{Type: "open-meteo"},
        Cache: config.Cache{
            Enabled: true,
            TTL:     time.Minute,
        },
    }

    provider := getProvider(cfg, testLogger{})

    if _, ok := provider.(*weather.WeatherInfo); !ok {
        t.Fatalf("expected *weather.WeatherInfo, got %T", provider)
    }
}

func TestGetProviderPogoda(t *testing.T) {
    cfg := config.Config{
        Provider: config.Provider{Type: "pogoda"},
    }

    provider := getProvider(cfg, testLogger{})

    if _, ok := provider.(*pogodaby.Pogoda); !ok {
        t.Fatalf("expected *pogodaby.Pogoda, got %T", provider)
    }
}

func TestGetProviderFallback(t *testing.T) {
    cfg := config.Config{
        Provider: config.Provider{Type: "unknown"},
        Cache: config.Cache{
            Enabled: true,
            TTL:     time.Minute,
        },
    }

    provider := getProvider(cfg, testLogger{})

    if _, ok := provider.(*weather.WeatherInfo); !ok {
        t.Fatalf("expected fallback to *weather.WeatherInfo, got %T", provider)
    }
}

var _ cli.Logger = testLogger{}
