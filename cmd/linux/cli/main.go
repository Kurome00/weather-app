package main

import (
	"fmt"
	"os"

	"github.com/Kurome00/weather-app.git/internal/pkg/app/cli"
	"github.com/Kurome00/weather-app.git/internal/pkg/config"
	"github.com/Kurome00/weather-app.git/internal/pkg/flags"
	"github.com/Kurome00/weather-app.git/internal/pkg/providers"
)

func main() {
	arguments := flags.Parse()

	r, err := os.Open(arguments.Path)
	if err != nil {
		fmt.Printf("Ошибка открытия файла конфигурации: %v\n", err)
		os.Exit(1)
	}
	defer r.Close()

	cfg, err := config.Parse(r)
	if err != nil {
		fmt.Printf("Ошибка парсинга конфигурации: %v\n", err)
		os.Exit(1)
	}

	logger := cli.NewConsoleLogger(true)

	wi := providers.GetProvider(cfg, logger)

	app := cli.New(logger, wi, cfg)
	if err := app.Run(); err != nil {
		logger.Error(fmt.Sprintf("Ошибка приложения: %v", err))
		os.Exit(1)
	}
}
