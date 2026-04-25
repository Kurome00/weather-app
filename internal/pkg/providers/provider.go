package providers

import (
	pogodaby "github.com/Kurome00/weather-app.git/internal/adapters/pogoda_by"
	"github.com/Kurome00/weather-app.git/internal/adapters/weather"
	"github.com/Kurome00/weather-app.git/internal/pkg/app/cli"
	"github.com/Kurome00/weather-app.git/internal/pkg/config"
)

func GetProvider(c config.Config, l cli.Logger) cli.WeatherInfo {
	var wi cli.WeatherInfo

	switch c.P.Type {
	case "open-meteo":
		wi = weather.New(l)
	case "pogoda":
		wi = pogodaby.New(l)
	default:
		wi = weather.New(l)
	}

	return wi
}
