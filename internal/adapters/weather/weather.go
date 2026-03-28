package weather

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
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

type WeatherInfo struct {
    logger   Logger
    current  current
    isLoaded bool
}

func New(logger Logger) *WeatherInfo {
    return &WeatherInfo{
        logger: logger,
    }
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
    wi.isLoaded = true
    return nil
}

func (wi *WeatherInfo) GetTemperature(lat, long float64) (float32, error) {
    if !wi.isLoaded {
        if err := wi.getWeatherInfo(lat, long); err != nil {
            return 0, err
        }
    }
    return wi.current.Temp, nil
}