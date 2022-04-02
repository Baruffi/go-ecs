package impl

import (
	"example.com/v0/src/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func setup() ecs.Ecs {
	cfg := pixelgl.WindowConfig{
		Title:  "Proto countries",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	registry := ecs.NewRegistry()
	testScene := ecs.NewScene(registry, ecs.UpdaterFunc(updateFromRegistry))
	setupScene(testScene)

	renderer := PixelRenderer{
		win:        win,
		clearColor: colornames.Black,
		drawers:    []ecs.Drawer{ecs.DrawerFunc(drawFromRegistry(registry))},
	}

	clock := ecs.Clock{}

	game := ecs.Ecs{
		Clock:    clock,
		Renderer: renderer,
		Scene:    testScene,
	}

	return game
}

func Run() {
	game := setup()

	game.Loop()
}
