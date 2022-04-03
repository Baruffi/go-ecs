package exampleScene

import (
	"fmt"
	"time"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

func New(win *pixelgl.Window) *ecs.Scene {
	registry := ecs.NewRegistry()
	updater := &EntityUpdater{}
	testScene := ecs.NewScene(registry, updater)
	configureScene(testScene, updater, win)

	return testScene
}

func configureScene(s *ecs.Scene, updater *EntityUpdater, win *pixelgl.Window) {
	timeComponent := &components.TimeComponent{
		Ticker: time.NewTicker(time.Second),
		Format: "Mon, 02 Jan 2006 15:04:05 MST",
	}
	cameraComponent := &components.CameraComponent{
		CamPos:       pixel.ZV,
		CamSpeed:     500.0,
		CamZoom:      1.0,
		CamZoomSpeed: 1.2,
		Active:       true,
	}
	textComponent := &components.TextComponent{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textComponent.Init(pixel.V(10, 10), atlas)

	worldMap := s.CreateEntity()

	ecs.AddComponent(worldMap, timeComponent)
	ecs.AddComponent(worldMap, cameraComponent)
	ecs.AddComponent(worldMap, textComponent)

	ecs.AddComponentGroup[components.Drawable](worldMap, textComponent)

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
		if cameraComponent, ok := ecs.GetComponent[*components.CameraComponent](e); ok {
			if cameraComponent.Active {
				if u.window.Pressed(pixelgl.MouseButton1) {
					mouseDelta := u.window.MousePosition().Sub(u.window.MousePreviousPosition())
					cameraComponent.Move(mouseDelta)
				}
				cameraComponent.Scroll(u.window.MouseScroll())
				cameraComponent.Update()

				u.window.SetMatrix(cameraComponent.Cam)
			}
		}
		if textComponent, ok := ecs.GetComponent[*components.TextComponent](e); ok {
			if timeComponent, ok := ecs.GetComponent[*components.TimeComponent](e); ok {
				select {
				case <-timeComponent.Ticker.C:
					textComponent.Clear()
					timeComponent.Time = time.Now()
					timeStr := fmt.Sprintf("TIME: %s", timeComponent.String())
					textComponent.Write(timeStr)
				default:
				}
			}
		}
	}
}
