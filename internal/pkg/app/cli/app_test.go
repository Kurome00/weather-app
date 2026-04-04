package cli

import (
    "testing"
)

type mockWeatherInfo struct {
    getTemperatureFunc func(lat, long float64) (float32, error)
    clearCacheCalled   bool
}

func (m *mockWeatherInfo) GetTemperature(lat, long float64) (float32, error) {
    if m.getTemperatureFunc != nil {
        return m.getTemperatureFunc(lat, long)
    }
    return 25.5, nil
}

func (m *mockWeatherInfo) ClearCache() {
    m.clearCacheCalled = true
}

type mockLoggerTest struct {
    infoCalled  bool
    debugCalled bool
    errorCalled bool
}

func (m *mockLoggerTest) Info(msg string) {
    m.infoCalled = true
}

func (m *mockLoggerTest) Debug(msg string) {
    m.debugCalled = true
}

func (m *mockLoggerTest) Error(msg string) {
    m.errorCalled = true
}

func TestNewCliApp(t *testing.T) {
    logger := &mockLoggerTest{}
    wi := &mockWeatherInfo{}

    app := New(logger, wi, 53.6688, 23.8223)

    if app == nil {
        t.Error("New() returned nil")
    }

    if app.logger != logger {
        t.Error("Logger not set correctly")
    }

    if app.wi != wi {
        t.Error("WeatherInfo not set correctly")
    }
}

func TestCliAppRunSuccess(t *testing.T) {
    logger := &mockLoggerTest{}
    wi := &mockWeatherInfo{
        getTemperatureFunc: func(lat, long float64) (float32, error) {
            return 25.5, nil
        },
    }

    app := New(logger, wi, 53.6688, 23.8223)
    err := app.Run()

    if err != nil {
        t.Errorf("Run() returned error: %v", err)
    }

    if !logger.infoCalled {
        t.Error("Info should be called")
    }
}