package impl

import (
	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/scenes/mainScene"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	MainSceneId = "mainScene"
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
	scenes := map[string]*ecs.Scene{
		MainSceneId: mainScene,
	}
	stage := ecs.NewStage(MainSceneId, scenes)

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

		for _, group := range ecs.MapGroup[components.Drawer](scene) {
			for _, drawer := range group {
				drawer.Draw(win)
			}
		}

		win.Update()
	}
}
