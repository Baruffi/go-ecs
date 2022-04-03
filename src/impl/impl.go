package impl

import (
	"example.com/v0/src/engine"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func setup() engine.Engine[*PixelRenderer] {
	cfg := pixelgl.WindowConfig{
		Title:  "Proto countries",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	running := func() bool {
		return !win.Closed()
	}

	clock := engine.Clock{}

	render := engine.Render[*PixelRenderer]{
		Renderer: &PixelRenderer{
			window:     win,
			clearColor: colornames.Black,
		},
		DrawHandler: engine.DrawHandlerFunc[*PixelRenderer](DrawTest),
	}

	testScene := SetupScene(win)

	engine := engine.NewEngine(engine.StatefulRunnerFunc(running), clock, render, testScene)

	return engine
}

func Run() {
	engine := setup()

	engine.Run()
}
