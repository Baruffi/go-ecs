package impl

import (
	"fmt"
	"time"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/scenes/mainScene"
	"example.com/v0/src/impl/systems"
	"example.com/v0/src/impl/tools"
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

func setupSystems(win *pixelgl.Window) (*systems.EventSystem, *systems.DrawSystem) {
	eventSystem := systems.NewEventSystem(10, 10, 5*time.Second)
	drawSystem := systems.NewDrawSystem(win)

	return eventSystem, drawSystem
}

func setupStage(win *pixelgl.Window, eventSystem *systems.EventSystem, drawSystem *systems.DrawSystem) *ecs.Stage {
	scenes := map[string]*ecs.Scene{
		MainSceneId: mainScene.NewScene(win, eventSystem, drawSystem),
	}
	stage := ecs.NewStage(MainSceneId, scenes)

	return stage
}

func setupClock() *tools.Clock {
	clock := &tools.Clock{}
	clock.Init(-1)

	return clock
}

func Run() {
	win := setupWindow()
	eventSystem, drawSystem := setupSystems(win)
	stage := setupStage(win, eventSystem, drawSystem)
	clock := setupClock()

	var (
		frames = 0
		second = time.NewTicker(time.Second)
	)
	for !win.Closed() {
		dt := clock.Dt()

		scene := stage.GetScene()
		scene.Update(dt)

		eventSystem.Execute()

		win.Clear(colornames.Black)

		drawSystem.Execute()

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
