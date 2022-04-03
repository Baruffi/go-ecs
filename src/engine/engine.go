package engine

import "example.com/v0/src/engine/ecs"

type Engine[R RendererType] struct {
	StatefulRunner
	clock       Clock
	render      Render[R]
	activeScene *ecs.Scene
}

func NewEngine[R RendererType](n StatefulRunner, c Clock, r Render[R], s *ecs.Scene) Engine[R] {
	return Engine[R]{
		StatefulRunner: n,
		clock:          c,
		render:         r,
		activeScene:    s,
	}
}

func (e *Engine[R]) SetScene(s *ecs.Scene) {
	e.activeScene = s
}

func (e *Engine[R]) Init() {
	e.clock.Init()
}

func (e *Engine[R]) Run() {
	e.Init()
	for e.Running() {
		e.Step()
	}
}

func (e *Engine[R]) Step() {
	e.clock.Tick()

	e.activeScene.Update(e.clock.dt)

	e.render.Clear()
	e.render.BeginScene(e.activeScene)
	e.render.DrawScene()
	e.render.EndScene()
}
