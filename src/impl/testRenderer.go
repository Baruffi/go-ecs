package impl

import (
	"image/color"

	"example.com/v0/src/engine"
	"example.com/v0/src/engine/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Drawable interface {
	Draw(pixel.Target)
}
type PixelRenderer struct {
	window     *pixelgl.Window
	clearColor color.RGBA
	emitters   map[ecs.EntityId]*RenderEventEmitterComponent
}

func (r *PixelRenderer) Consume(events engine.Events[engine.RenderCall]) {
	for call, data := range events {
		switch call {
		case engine.DrawCall:
			for _, v := range data {
				if drawable, ok := v.(Drawable); ok {
					drawable.Draw(r.window)
				}
			}
		case engine.ClearCall:
			r.window.Clear(r.clearColor)
		}
	}
}

func (r *PixelRenderer) BeginScene(s *ecs.Scene) {
	for _, camera := range ecs.Map[*CameraComponent](s) {
		if camera.active {
			r.window.SetMatrix(camera.cam)
			break
		}
	}

	// Hmmmm
	r.emitters = ecs.Map[*RenderEventEmitterComponent](s)
}

func (r *PixelRenderer) EndScene() {
	r.window.Update()
}

func DrawTest(r *PixelRenderer) {
	for _, emitter := range r.emitters {
		r.Consume(emitter.Emit())
	}
}
