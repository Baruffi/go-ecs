package ecs

type Ecs struct {
	Clock
	Renderer
	*Scene
}

func (e *Ecs) Loop() {
	e.Init()
	for e.CanRender() {
		e.Tick()
		e.Clear()
		e.Update(e.dt)
		e.BeginScene()
		e.EndScene()
	}
}
