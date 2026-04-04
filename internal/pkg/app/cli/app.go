package cli

import (
    "fmt"
)

type cliApp struct {
    logger Logger
    wi     WeatherInfo
    lat    float64
    long   float64
}

type WeatherInfo interface {
    GetTemperature(float64, float64) (float32, error)
    ClearCache()
}

func New(logger Logger, wi WeatherInfo, lat, long float64) *cliApp {
    return &cliApp{
        logger: logger,
        wi:     wi,
        lat:    lat,
        long:   long,
    }
}

func (c *cliApp) Run() error {
    c.logger.Info("Запуск приложения для получения погоды")
    c.logger.Debug(fmt.Sprintf("Конфигурация: широта=%.4f, долгота=%.4f",
        c.lat, c.long))

    temp, err := c.wi.GetTemperature(c.lat, c.long)
    if err != nil {
        c.logger.Error(fmt.Sprintf("Ошибка получения температуры: %v", err))
        return err
    }

    result := fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", temp)
    c.logger.Info(result)
    fmt.Println(result)

    return nil
}