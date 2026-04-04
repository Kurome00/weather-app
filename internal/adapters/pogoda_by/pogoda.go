package pogodaby

import (
    "encoding/json"
    "fmt"
    "net/http"
)

const url = "https://pogoda.by/api/v2/weather-fact?station=26820"

type Logger interface {
    Info(string)
    Debug(string)
    Error(string)
}

type response struct {
    Temp float32 `json:"t"`
}

type Pogoda struct {
    logger Logger
}

func New(logger Logger) *Pogoda {
    return &Pogoda{logger: logger}
}

func (p *Pogoda) GetTemperature(_, _ float64) (float32, error) {
    resp, err := http.Get(url)
    if err != nil {
        p.logger.Error(fmt.Sprintf("не удалось получить данные от pogoda.by: %v", err))
        return 0, fmt.Errorf("не удалось получить данные от pogoda.by: %w", err)
    }
    defer func() {
        if err := resp.Body.Close(); err != nil {
            p.logger.Error(fmt.Sprintf("не удалось закрыть тело ответа pogoda.by: %v", err))
        }
    }()

    var data response
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        p.logger.Error(fmt.Sprintf("не удалось декодировать JSON pogoda.by: %v", err))
        return 0, fmt.Errorf("не удалось декодировать JSON pogoda.by: %w", err)
    }

    return data.Temp, nil
}

func (p *Pogoda) ClearCache() {}
