package weather

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"
)

const apiURL = "https://api.open-meteo.com/v1/forecast"

type Logger interface {
    Info(string)
    Debug(string)
    Error(string)
}

type current struct {
    Temp float32 `json:"temperature_2m"`
}

type response struct {
    Curr current `json:"current"`
}

type cacheEntry struct {
    temperature float32
    timestamp   time.Time
}

type WeatherInfo struct {
    logger       Logger
    current      current
    cache        map[string]cacheEntry
    mu           sync.RWMutex
    cacheTTL     time.Duration
    cacheEnabled bool
}

type CacheConfig struct {
    Enabled bool
    TTL     time.Duration
}

func New(logger Logger, cacheEnabled bool, cacheTTL time.Duration) *WeatherInfo {
    wi := &WeatherInfo{
        logger:       logger,
        cache:        make(map[string]cacheEntry),
        cacheTTL:     cacheTTL,
        cacheEnabled: cacheEnabled,
    }

    if wi.cacheEnabled {
        wi.logger.Debug(fmt.Sprintf("Кэширование включено, TTL: %v", wi.cacheTTL))
    } else {
        wi.logger.Debug("Кэширование отключено")
    }

    return wi
}

func (wi *WeatherInfo) getCacheKey(lat, long float64) string {
    return fmt.Sprintf("%f,%f", lat, long)
}

func (wi *WeatherInfo) getFromCache(lat, long float64) (float32, bool) {
    if !wi.cacheEnabled {
        return 0, false
    }

    wi.mu.RLock()
    defer wi.mu.RUnlock()

    key := wi.getCacheKey(lat, long)
    entry, exists := wi.cache[key]
    if !exists {
        return 0, false
    }

    if time.Since(entry.timestamp) > wi.cacheTTL {
        wi.logger.Debug(fmt.Sprintf("Кэш устарел для координат %s", key))
        return 0, false
    }

    wi.logger.Debug(fmt.Sprintf("Данные получены из кэша для %s", key))
    return entry.temperature, true
}

func (wi *WeatherInfo) saveToCache(lat, long float64, temp float32) {
    if !wi.cacheEnabled {
        return
    }

    wi.mu.Lock()
    defer wi.mu.Unlock()

    key := wi.getCacheKey(lat, long)
    wi.cache[key] = cacheEntry{
        temperature: temp,
        timestamp:   time.Now(),
    }
    wi.logger.Debug(fmt.Sprintf("Данные сохранены в кэш для %s", key))
}

func (wi *WeatherInfo) getWeatherInfo(lat, long float64) error {
    var respData response

    params := fmt.Sprintf(
        "latitude=%f&longitude=%f&current=temperature_2m",
        lat, long,
    )

    url := fmt.Sprintf("%s?%s", apiURL, params)
    wi.logger.Debug(fmt.Sprintf("URL успешно сгенерирован - %s", url))

    resp, err := http.Get(url)
    if err != nil {
        wi.logger.Error(fmt.Sprintf("не удалось получить данные о погоде: %v", err))
        return fmt.Errorf("не удалось получить данные о погоде из openmeteo: %w", err)
    }
    defer func() {
        if err := resp.Body.Close(); err != nil {
            wi.logger.Error(fmt.Sprintf("не удалось закрыть тело ответа: %v", err))
        }
    }()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        wi.logger.Error(fmt.Sprintf("не удалось прочитать данные из тела: %v", err))
        return fmt.Errorf("не удалось прочитать данные из ответа: %w", err)
    }

    wi.logger.Debug(fmt.Sprintf("данные успешно прочитаны, размер - %d байт", len(data)))

    if err := json.Unmarshal(data, &respData); err != nil {
        wi.logger.Error(fmt.Sprintf("не удалось распарсить JSON: %v", err))
        return fmt.Errorf("не удалось распарсить данные из ответа: %w", err)
    }

    wi.current = respData.Curr
    return nil
}

func (wi *WeatherInfo) GetTemperature(lat, long float64) (float32, error) {
    cachedTemp, found := wi.getFromCache(lat, long)
    if found {
        return cachedTemp, nil
    }

    if err := wi.getWeatherInfo(lat, long); err != nil {
        return 0, err
    }

    wi.saveToCache(lat, long, wi.current.Temp)
    return wi.current.Temp, nil
}

func (wi *WeatherInfo) ClearCache() {
    wi.mu.Lock()
    defer wi.mu.Unlock()
    wi.cache = make(map[string]cacheEntry)
    wi.logger.Info("Кэш очищен")
}