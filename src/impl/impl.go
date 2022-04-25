package impl

import (
	"fmt"
	"log"
	"os"
	"time"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/factories/countryFactory"
	"example.com/v0/src/impl/scenes/mainScene"
	"example.com/v0/src/impl/systems"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	MainSceneId = "mainScene"
)

func setupLog() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

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

func setupMainScene(win *pixelgl.Window, eventSystem *systems.EventSystem, drawSystem *systems.DrawSystem) (string, *ecs.Scene, ecs.Updater[*ecs.Stage]) {
	newMainScene := mainScene.NewScene()
	player, world, ui := mainScene.ConfigureScene(newMainScene, win, eventSystem, drawSystem)
	countryFactory := countryFactory.NewFactory(newMainScene, 0, win.Bounds().Center(), pixel.ZV, "EST", eventSystem, drawSystem)
	mainUpdater := mainScene.NewUpdater(countryFactory, win, player, world, ui)

	return MainSceneId, newMainScene, mainUpdater
}

func setupStage(win *pixelgl.Window, eventSystem *systems.EventSystem, drawSystem *systems.DrawSystem) *ecs.Stage {
	stage := ecs.NewStage()
	stage.Include(setupMainScene(win, eventSystem, drawSystem))
	stage.Start(MainSceneId)

	return stage
}

func Run() {
	setupLog()
	win := setupWindow()
	eventSystem, drawSystem := setupSystems(win)
	stage := setupStage(win, eventSystem, drawSystem)

	var (
		frames = 0
		second = time.NewTicker(time.Second)
	)
	for stage.IsActive() {
		stagehand := stage.GetStagehand()
		stagehand.Update()

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

		if win.Closed() {
			stage.End()
		}
	}
}
