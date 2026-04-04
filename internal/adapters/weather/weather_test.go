package weather

import (
    "testing"
    "time"
)

type mockLogger struct {
    debugCalled bool
    infoCalled  bool
    errorCalled bool
    lastMessage string
}

func (m *mockLogger) Info(msg string) {
    m.infoCalled = true
    m.lastMessage = msg
}

func (m *mockLogger) Debug(msg string) {
    m.debugCalled = true
    m.lastMessage = msg
}

func (m *mockLogger) Error(msg string) {
    m.errorCalled = true
    m.lastMessage = msg
}

func TestNewWeatherInfo(t *testing.T) {
    logger := &mockLogger{}

    wi := New(logger, true, 5*time.Minute)

    if wi == nil {
        t.Error("New() returned nil")
    }

    if wi.cacheEnabled != true {
        t.Error("Cache should be enabled")
    }

    if wi.cacheTTL != 5*time.Minute {
        t.Errorf("Expected TTL 5m, got %v", wi.cacheTTL)
    }
}

func TestNewWeatherInfoCacheDisabled(t *testing.T) {
    logger := &mockLogger{}

    wi := New(logger, false, 5*time.Minute)

    if wi.cacheEnabled != false {
        t.Error("Cache should be disabled")
    }
}

func TestGetCacheKey(t *testing.T) {
    logger := &mockLogger{}

    wi := New(logger, true, 5*time.Minute)

    key1 := wi.getCacheKey(53.6688, 23.8223)
    key2 := wi.getCacheKey(53.6688, 23.8223)
    key3 := wi.getCacheKey(55.7558, 37.6173)

    if key1 != key2 {
        t.Error("Same coordinates should produce same cache key")
    }

    if key1 == key3 {
        t.Error("Different coordinates should produce different cache key")
    }
}

func TestSaveAndGetFromCache(t *testing.T) {
    logger := &mockLogger{}

    wi := New(logger, true, 5*time.Minute)

    wi.saveToCache(53.6688, 23.8223, 25.5)

    temp, found := wi.getFromCache(53.6688, 23.8223)

    if !found {
        t.Error("Cache entry not found")
    }

    if temp != 25.5 {
        t.Errorf("Expected 25.5, got %f", temp)
    }
}

func TestGetFromCacheNotFound(t *testing.T) {
    logger := &mockLogger{}

    wi := New(logger, true, 5*time.Minute)

    _, found := wi.getFromCache(99.9999, 99.9999)

    if found {
        t.Error("Should not find cache entry for non-existent key")
    }
}

func TestGetFromCacheExpired(t *testing.T) {
    logger := &mockLogger{}

    wi := New(logger, true, 1*time.Nanosecond)

    wi.saveToCache(53.6688, 23.8223, 25.5)

    time.Sleep(1 * time.Millisecond)

    _, found := wi.getFromCache(53.6688, 23.8223)

    if found {
        t.Error("Cache entry should be expired")
    }
}

func TestGetFromCacheDisabled(t *testing.T) {
    logger := &mockLogger{}

    wi := New(logger, false, 5*time.Minute)

    wi.saveToCache(53.6688, 23.8223, 25.5)

    _, found := wi.getFromCache(53.6688, 23.8223)

    if found {
        t.Error("Should not get from cache when cache is disabled")
    }
}

func TestClearCache(t *testing.T) {
    logger := &mockLogger{}

    wi := New(logger, true, 5*time.Minute)

    wi.saveToCache(53.6688, 23.8223, 25.5)
    wi.saveToCache(55.7558, 37.6173, 15.5)

    wi.ClearCache()

    if len(wi.cache) != 0 {
        t.Errorf("Cache should be empty, but has %d entries", len(wi.cache))
    }
}