package guisettings

type TextWidget interface {
	Render() any
	SetText(text string)
}

type Window interface {
	Resize(ws WindowSize) error
	UpdateTemperature(t float32) error
	SetTemperatureWidget(tw TextWidget) error
	Render() error
}

type AppRunner interface {
	Run()
}

type Provider interface {
	CreateWindow(name string, size WindowSize) (Window, error)
	GetAppRunner() AppRunner
	GetTextWidget(text string) TextWidget
}
