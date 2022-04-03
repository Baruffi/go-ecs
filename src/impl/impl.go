package impl

import (
	"example.com/v0/src/engine"
	"example.com/v0/src/engine/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func setup() (*pixelgl.Window, ecs.Stage, engine.Clock) {
	cfg := pixelgl.WindowConfig{
		Title:  "Proto countries",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	stage := SetupStage(win)

	clock := engine.Clock{}

	return win, stage, clock
}

func Run() {
	win, stage, clock := setup()

	clock.Init()
	for !win.Closed() {
		dt := clock.Tick()

		scene := stage.GetScene()
		scene.Update(dt)

		win.Clear(colornames.Black)

		var drawables []Drawable2
		for _, group := range ecs.MapGroup[Drawable2](scene) {
			for _, drawable := range group {
				drawables = append(drawables, drawable)
			}
		}
		drawer := Drawer{
			target:    win,
			drawables: drawables,
		}
		drawer.Draw()

		win.Update()
	}
}
