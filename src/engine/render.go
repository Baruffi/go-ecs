package engine

import (
	"example.com/v0/src/engine/ecs"
)

type RenderCall uint

const (
	ClearCall RenderCall = iota
	DrawCall
)

type Renderer interface {
	BeginScene(*ecs.Scene)
	Consume(Events[RenderCall])
	EndScene()
}

type RendererType interface{}

type DrawHandler[R RendererType] interface {
	Draw(R)
}

type DrawHandlerFunc[R RendererType] func(R)

func (f DrawHandlerFunc[R]) Draw(r R) {
	f(r)
}

type Render[R RendererType] struct {
	Renderer
	DrawHandler[R]
}

// DrawScene - Call the draw handler on the current renderer
func (r *Render[R]) DrawScene() {
	r.Draw(r.Renderer.(R))
}

// Clear - Method to simplify clear calls
func (r *Render[R]) Clear() {
	events := make(Events[RenderCall])
	events[ClearCall] = nil

	r.Consume(events)
}
