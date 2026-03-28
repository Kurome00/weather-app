package cli

import (
    "fmt"
)

type Config struct {
    Latitude  float64
    Longitude float64
    DebugMode bool
}

type cliApp struct {
    logger Logger
    wi     WeatherInfo
    config Config
}

type WeatherInfo interface {
    GetTemperature(float64, float64) (float32, error)
}

func New(logger Logger, wi WeatherInfo, config Config) *cliApp {
    return &cliApp{
        logger: logger,
        wi:     wi,
        config: config,
    }
}

func (c *cliApp) Run() error {
    c.logger.Info("Запуск приложения для получения погоды")
    c.logger.Debug(fmt.Sprintf("Конфигурация: широта=%.4f, долгота=%.4f",
        c.config.Latitude, c.config.Longitude))

    temp, err := c.wi.GetTemperature(c.config.Latitude, c.config.Longitude)
    if err != nil {
        c.logger.Error(fmt.Sprintf("Ошибка получения температуры: %v", err))
        return err
    }

    result := fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", temp)
    c.logger.Info(result)
    fmt.Println(result)

    return nil
}