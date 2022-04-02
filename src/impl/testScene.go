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

func testUpdater(r *ecs.Registry, dt float64) {
	for _, e := range ecs.View[*TimeComponent](r) {
		timeData, _ := ecs.Get[*TimeComponent](r, e)
		if textData, hasText := ecs.Get[*TextComponent](r, e); hasText {
			timeData.Update()

			select {
			case <-timeData.ticker.C:
				textData.Write(fmt.Sprintf("FPS: %d", timeData.frames))
				timeData.frames = 0
			default:
			}
		}
	}
}

func testDrawer(r *ecs.Registry, t pixel.Target) {
	for _, e := range ecs.View[*TextComponent](r) {
		drawable, _ := ecs.Get[*TextComponent](r, e)
		drawable.Draw(t)
	}
}

func setup() (*pixelgl.Window, *ecs.Registry) {
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

	timeData := &TimeComponent{
		ticker: time.NewTicker(time.Second),
	}
	cameraData := &CameraComponent{
		camPos:       pixel.ZV,
		camSpeed:     500.0,
		camZoom:      1.0,
		camZoomSpeed: 1.2,
	}
	textData := &TextComponent{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textData.Init(pixel.V(100, 500), atlas)

	worldMap := ecs.NewEntity()

	ecs.Link(registry, worldMap, timeData)
	ecs.Link(registry, worldMap, textData)
	ecs.Link(registry, worldMap, cameraData)

	return win, registry
}

func Run() {
	win, registry := setup()
	testScene := ecs.Scene{
		Registry: registry,
		Updater:  ecs.UpdaterFunc(testUpdater),
		Drawer:   ecs.DrawerFunc(testDrawer),
	}

	updateCh := time.NewTicker(time.Second)
	frames := 0
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		testScene.Update(dt)

		win.Clear(colornames.Black)
		testScene.Draw(win)
		win.Update()

		frames++
		select {
		case <-updateCh.C:
			win.SetTitle(fmt.Sprintf("Proto countries | FPS: %d", frames))
			frames = 0
		default:
		}
	}
}
