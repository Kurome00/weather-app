package fyne

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	guisettings "github.com/Kurome00/weather-app.git/internal/domain/gui_settings"
)

type window struct {
	w  fyne.Window
	tw guisettings.TextWidget
}

func NewW(win fyne.Window) *window {
	return &window{w: win}
}

func (w *window) Resize(ws guisettings.WindowSize) error {
	if ws.IsFull() {
		w.w.SetFullScreen(true)
	} else {
		w.w.Resize(fyne.Size{Width: float32(ws.Width()), Height: float32(ws.Heigth())})
	}
	return nil
}

func (w *window) UpdateTemperature(t float32) error {
	w.tw.SetText(fmt.Sprintf("За окном - %f", t))
	return nil
}

func (w *window) SetTemperatureWidget(tw guisettings.TextWidget) error {
	w.tw = tw
	label := tw.Render()
	labelFyne := label.(*widget.Label)
	center := container.NewCenter(labelFyne)
	w.w.SetContent(center)
	return nil
}

func (w *window) Render() error {
	w.w.Show()
	return nil
}
