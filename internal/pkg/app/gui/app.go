package gui

import (
	"fmt"

	guisettings "github.com/Kurome00/weather-app.git/internal/domain/gui_settings"
	"github.com/Kurome00/weather-app.git/internal/pkg/app/cli"
	"github.com/Kurome00/weather-app.git/internal/pkg/config"
)

type App struct {
	logger   cli.Logger
	provider guisettings.Provider
	wi       cli.WeatherInfo
	config   config.Config
}

func New(logger cli.Logger, provider guisettings.Provider, wi cli.WeatherInfo, cfg config.Config) *App {
	return &App{logger: logger, provider: provider, wi: wi, config: cfg}
}

func (a *App) Run() error {
	a.logger.Info("Запуск GUI приложения для получения погоды")

	w, err := a.provider.CreateWindow("Информер погоды", guisettings.NewWS(480, 320))
	if err != nil {
		return fmt.Errorf("не удалось создать окно: %w", err)
	}

	tw := a.provider.GetTextWidget("Загрузка данных о погоде...")
	if err := w.SetTemperatureWidget(tw); err != nil {
		return fmt.Errorf("не удалось установить виджет температуры: %w", err)
	}

	temp, err := a.wi.GetTemperature(a.config.L.Lat, a.config.L.Long)
	if err != nil {
		return fmt.Errorf("не удалось получить температуру: %w", err)
	}

	if err := w.UpdateTemperature(temp); err != nil {
		return fmt.Errorf("не удалось обновить температуру: %w", err)
	}

	if err := w.Render(); err != nil {
		return fmt.Errorf("не удалось отрисовать окно: %w", err)
	}

	a.provider.GetAppRunner().Run()
	return nil
}
