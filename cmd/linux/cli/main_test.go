package main

import (
	"testing"

	pogodaby "github.com/Kurome00/weather-app.git/internal/adapters/pogoda_by"
	"github.com/Kurome00/weather-app.git/internal/adapters/weather"
	"github.com/Kurome00/weather-app.git/internal/pkg/app/cli"
	"github.com/Kurome00/weather-app.git/internal/pkg/config"
	"github.com/Kurome00/weather-app.git/internal/pkg/providers"
)

type testLogger struct{}

func (testLogger) Info(string)  {}
func (testLogger) Debug(string) {}
func (testLogger) Error(string) {}

func TestGetProviderOpenMeteo(t *testing.T) {
	cfg := config.Config{P: config.Provider{Type: "open-meteo"}}

	provider := providers.GetProvider(cfg, testLogger{})

	if _, ok := provider.(*weather.WeatherInfo); !ok {
		t.Fatalf("expected *weather.WeatherInfo, got %T", provider)
	}
}

func TestGetProviderPogoda(t *testing.T) {
	cfg := config.Config{P: config.Provider{Type: "pogoda"}}

	provider := providers.GetProvider(cfg, testLogger{})

	if _, ok := provider.(*pogodaby.Pogoda); !ok {
		t.Fatalf("expected *pogodaby.Pogoda, got %T", provider)
	}
}

func TestGetProviderFallback(t *testing.T) {
	cfg := config.Config{P: config.Provider{Type: "unknown"}}

	provider := providers.GetProvider(cfg, testLogger{})

	if _, ok := provider.(*weather.WeatherInfo); !ok {
		t.Fatalf("expected fallback to *weather.WeatherInfo, got %T", provider)
	}
}

var _ cli.Logger = testLogger{}
