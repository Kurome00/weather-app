package main

import (
	"os"

	"github.com/Kurome00/weather-app.git/internal/pkg/app/cli"
	"github.com/Kurome00/weather-app.git/internal/pkg/app/gui"
	"github.com/Kurome00/weather-app.git/internal/pkg/config"
	"github.com/Kurome00/weather-app.git/internal/pkg/flags"
	fygui "github.com/Kurome00/weather-app.git/internal/pkg/gui/fyne"
	"github.com/Kurome00/weather-app.git/internal/pkg/providers"
)

func main() {
	arguments := flags.Parse()

	r, err := os.Open(arguments.Path)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	c, err := config.Parse(r)
	if err != nil {
		panic(err)
	}

	l := cli.NewConsoleLogger(true)
	provider := providers.GetProvider(c, l)
	p := fygui.NewP()
	g := gui.New(l, p, provider, c)
	if err := g.Run(); err != nil {
		panic(err)
	}
}
