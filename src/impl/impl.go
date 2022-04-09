package impl

import (
	"fmt"
	"time"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/managers"
	"example.com/v0/src/impl/scenes/mainScene"
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

func setupManagers(win *pixelgl.Window) (managers.EventManager, managers.DrawerManager) {
	eventManager := managers.NewEventManager()
	drawerManager := managers.NewDrawerManager(win)

	return eventManager, drawerManager
}

func setupStage(win *pixelgl.Window, drawerManager managers.DrawerManager) ecs.Stage {
	mainScene := mainScene.NewScene(win, drawerManager)
	scenes := map[string]*ecs.Scene{
		MainSceneId: mainScene,
	}
	stage := ecs.NewStage(MainSceneId, scenes)

	return stage
}

func setupClock() tools.Clock {
	clock := tools.Clock{}
	clock.Init(-1)

	return clock
}

func Run() {
	win := setupWindow()
	_, drawerManager := setupManagers(win)
	stage := setupStage(win, drawerManager)
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

		drawerManager.Execute()

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
