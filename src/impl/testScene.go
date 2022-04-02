package impl

import (
	"fmt"
	"time"

	"example.com/v0/src/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func setupScene(s *ecs.Scene) {
	timeComponent := &TimeComponent{
		ticker: time.NewTicker(time.Second),
	}
	cameraComponent := &CameraComponent{
		camPos:       pixel.ZV,
		camSpeed:     500.0,
		camZoom:      1.0,
		camZoomSpeed: 1.2,
	}
	textComponent := &TextComponent{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textComponent.Init(pixel.V(100, 500), atlas)

	worldMap := s.CreateEntity()

	ecs.AddComponent(worldMap, timeComponent)
	ecs.AddComponent(worldMap, textComponent)
	ecs.AddComponent(worldMap, cameraComponent)
}

func updateFromRegistry(r *ecs.Registry, dt float64) {
	for _, e := range ecs.View[*TimeComponent](r) {
		timeComponent, _ := ecs.Get[*TimeComponent](r, e)
		if textComponent, hasText := ecs.Get[*TextComponent](r, e); hasText {
			timeComponent.Update()

			select {
			case <-timeComponent.ticker.C:
				textComponent.Clear()
				textComponent.Write(fmt.Sprintf("FPS: %d", timeComponent.frames))
				timeComponent.frames = 0
			default:
			}
		}
	}
}

func drawFromRegistry(r *ecs.Registry) func(s ecs.Surface) {
	testDrawer := func(s ecs.Surface) {
		t := s.(pixel.Target)
		for _, e := range ecs.View[*TextComponent](r) {
			textComponent, _ := ecs.Get[*TextComponent](r, e)
			textComponent.Draw(t)
		}
	}

	return testDrawer
}

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
