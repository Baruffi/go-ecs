package impl

import (
	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/scenes/mainScene"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func setupWindowAndClock() (*pixelgl.Window, Clock) {
	cfg := pixelgl.WindowConfig{
		Title:  "Proto countries",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	clock := Clock{}

	return win, clock
}

func setupStage(win *pixelgl.Window) ecs.Stage {
	mainScene := mainScene.NewScene(win)
	scenes := []*ecs.Scene{mainScene}
	stage := ecs.NewStage(0, scenes)

	return stage
}

func Run() {
	win, clock := setupWindowAndClock()
	stage := setupStage(win)

	clock.Init()
	for !win.Closed() {
		dt := clock.Tick()

		scene := stage.GetScene()
		scene.Update(dt)

		win.Clear(colornames.Black)

		for _, group := range ecs.MapGroup[components.Drawable](scene) {
			for _, drawable := range group {
				drawable.Draw(win)
			}
		}

		win.Update()
	}
}
