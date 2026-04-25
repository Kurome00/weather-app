package guisettings

type WindowSize struct {
	width  int
	heigth int
}

func NewWS(w, h int) WindowSize {
	return WindowSize{width: w, heigth: h}
}

func (ws WindowSize) IsFull() bool {
	return ws.width == 0 && ws.heigth == 0
}

func (ws WindowSize) Width() int {
	return ws.width
}

func (ws WindowSize) Heigth() int {
	return ws.heigth
}
