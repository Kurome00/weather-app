package cli

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type Config struct {
    Latitude  float64
    Longitude float64
    DebugMode bool
}

type cliApp struct {
    logger Logger  // Это использует интерфейс из logger.go
    config Config
}

func New(logger Logger, config Config) *cliApp {
    return &cliApp{
        logger: logger,
        config: config,
    }
}

func (c *cliApp) Run() error {
    c.logger.Info("Запуск приложения для получения погоды")
    c.logger.Debug(fmt.Sprintf("Конфигурация: широта=%.4f, долгота=%.4f", 
        c.config.Latitude, c.config.Longitude))
    
    type Current struct {
        Temp float32 `json:"temperature_2m"`
    }

    type Response struct {
        Curr Current `json:"current"`
    }

    var response Response

    params := fmt.Sprintf(
        "latitude=%f&longitude=%f&current=temperature_2m",
        c.config.Latitude,
        c.config.Longitude,
    )

    url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?%s", params)
    
    c.logger.Debug(fmt.Sprintf("Отправка запроса к API: %s", url))

    resp, err := http.Get(url)
    if err != nil {
        c.logger.Error(fmt.Sprintf("Ошибка HTTP запроса: %s", err.Error()))
        return fmt.Errorf("не удалось получить данные о погоде из openmeteo: %v", err)
    }
    defer func() {
        if err := resp.Body.Close(); err != nil {
            c.logger.Error(fmt.Sprintf("Ошибка при закрытии тела ответа: %s", err.Error()))
        }
    }()

    c.logger.Debug(fmt.Sprintf("Статус ответа: %s", resp.Status))

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        c.logger.Error(fmt.Sprintf("Ошибка чтения данных: %s", err.Error()))
        return fmt.Errorf("не удалось прочитать данные из ответа: %v", err)
    }

    c.logger.Debug(fmt.Sprintf("Получено %d байт данных", len(data)))

    if err := json.Unmarshal(data, &response); err != nil {
        c.logger.Error(fmt.Sprintf("Ошибка парсинга JSON: %s", err.Error()))
        return fmt.Errorf("не удалось обработать данные из ответа: %v", err)
    }

    result := fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", response.Curr.Temp)
    c.logger.Info(result)
    fmt.Println(result)

    return nil
}