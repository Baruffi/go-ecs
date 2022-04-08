package impl

import (
	"fmt"
	"sort"
	"time"

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

func setupWindow() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  "Proto countries",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return win
}

func setupStage(win *pixelgl.Window) ecs.Stage {
	mainScene := mainScene.NewScene(win)
	scenes := map[string]*ecs.Scene{
		MainSceneId: mainScene,
	}
	stage := ecs.NewStage(MainSceneId, scenes)

	return stage
}

func setupClock() Clock {
	clock := Clock{}
	clock.Init(-1)

	return clock
}

func Run() {
	win := setupWindow()
	stage := setupStage(win)
	clock := setupClock()

	var (
		frames = 0
		second = time.NewTicker(time.Second)
	)
	for !win.Closed() {
		dt := clock.Dt()

		scene := stage.GetScene()
		scene.Update(dt)

		win.Clear(colornames.Black)

		drawQueue := &components.SortableDrawerQueue{}
		drawQueue.Init()
		for _, group := range ecs.MapGroup[components.Drawer](scene) {
			for _, drawer := range group {
				drawQueue.Add(drawer)
			}
		}
		sort.Sort(drawQueue)
		for _, drawer := range drawQueue.Drawers {
			drawer.Draw(win)
		}

		win.Update()

		frames++
		select {
		case <-second.C:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", "Proto countries", frames))
			frames = 0
		default:
		}

		clock.Tick()
	}
}
