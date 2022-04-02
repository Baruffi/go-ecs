package impl

import (
	"image/color"

	"example.com/v0/src/ecs"
	"github.com/faiface/pixel/pixelgl"
)

type PixelRenderer struct {
	win        *pixelgl.Window
	clearColor color.RGBA
	drawers    []ecs.Drawer
}

func (r PixelRenderer) CanRender() bool {
	return !r.win.Closed()
}

func (r PixelRenderer) Clear() {
	r.win.Clear(r.clearColor)
}

func (r PixelRenderer) BeginScene() {
	for _, drawer := range r.drawers {
		drawer.Draw(r.win)
	}
}

func (r PixelRenderer) EndScene() {
	r.win.Update()
}
