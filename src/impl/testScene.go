package impl

import (
	"fmt"
	"time"

	"example.com/v0/src/engine"
	"example.com/v0/src/engine/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

func setupScene(s *ecs.Scene, updater *EntityUpdater, win *pixelgl.Window) {
	timeComponent := &TimeComponent{
		ticker: time.NewTicker(time.Second),
		format: "Mon, 02 Jan 2006 15:04:05 MST",
	}
	cameraComponent := &CameraComponent{
		camPos:       pixel.ZV,
		camSpeed:     500.0,
		camZoom:      1.0,
		camZoomSpeed: 1.2,
		active:       true,
	}
	textComponent := &TextComponent{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textComponent.Init(pixel.V(10, 10), atlas)
	renderEmitterComponent := &RenderEventEmitterComponent{}
	renderEmitterComponent.Clear()
	renderEmitterComponent.Add(engine.DrawCall, textComponent)

	worldMap := s.CreateEntity()

	ecs.AddComponent(worldMap, timeComponent)
	ecs.AddComponent(worldMap, cameraComponent)
	ecs.AddComponent(worldMap, textComponent)
	ecs.AddComponent(worldMap, renderEmitterComponent)

	// Hmmmm
	updater.window = win
	updater.entities = append(updater.entities, worldMap)
}

type EntityUpdater struct {
	window   *pixelgl.Window
	entities []ecs.Entity
}

func (u *EntityUpdater) Update(dt float64) {
	for _, e := range u.entities {
		if cameraComponent, ok := ecs.GetComponent[*CameraComponent](e); ok {
			if u.window.Pressed(pixelgl.MouseButton1) {
				mouseDelta := u.window.MousePosition().Sub(u.window.MousePreviousPosition())
				cameraComponent.Move(mouseDelta)
			}
			cameraComponent.Scroll(u.window.MouseScroll())
			cameraComponent.Update()
		}
		if textComponent, ok := ecs.GetComponent[*TextComponent](e); ok {
			if timeComponent, ok := ecs.GetComponent[*TimeComponent](e); ok {
				select {
				case <-timeComponent.ticker.C:
					textComponent.Clear()
					timeComponent.time = time.Now()
					timeStr := fmt.Sprintf("TIME: %s", timeComponent.Format())
					textComponent.Write(timeStr)
				default:
				}
			}
		}
	}
}

func SetupScene(win *pixelgl.Window) *ecs.Scene {
	registry := ecs.NewRegistry()
	updater := &EntityUpdater{}
	testScene := ecs.NewScene(registry, updater)
	setupScene(testScene, updater, win)

	return testScene
}
