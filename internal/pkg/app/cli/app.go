package cli

import (
    "fmt"

    "github.com/Kurome00/weather-app.git/internal/pkg/config"
)

type cliApp struct {
    logger Logger
    wi     WeatherInfo
    config config.Config
}

type WeatherInfo interface {
    GetTemperature(float64, float64) (float32, error)
    ClearCache()
}

func New(logger Logger, wi WeatherInfo, cfg config.Config) *cliApp {
    return &cliApp{
        logger: logger,
        wi:     wi,
        config: cfg,
    }
}

func (c *cliApp) Run() error {
    c.logger.Info("Запуск приложения для получения погоды")
    c.logger.Debug(fmt.Sprintf("Конфигурация: широта=%.4f, долгота=%.4f",
        c.config.L.Lat, c.config.L.Long))

    temp, err := c.wi.GetTemperature(c.config.L.Lat, c.config.L.Long)
    if err != nil {
        c.logger.Error(fmt.Sprintf("Ошибка получения температуры: %v", err))
        return err
    }

    result := fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", temp)
    c.logger.Info(result)
    fmt.Println(result)

    return nil
}