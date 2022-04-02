package impl

import (
	"fmt"
	"time"

	"example.com/v0/src/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

func setupScene(s *ecs.Scene) {
	timeComponent := &TimeComponent{
		ticker: time.NewTicker(time.Second),
		format: "Mon, 02 Jan 2006 15:04:05 MST",
	}
	cameraComponent := &CameraComponent{
		camPos:       pixel.ZV,
		camSpeed:     500.0,
		camZoom:      1.0,
		camZoomSpeed: 1.2,
	}
	textComponent := &TextComponent{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textComponent.Init(pixel.V(0, 0), atlas)

	worldMap := s.CreateEntity()

	ecs.AddComponent(worldMap, timeComponent)
	ecs.AddComponent(worldMap, textComponent)
	ecs.AddComponent(worldMap, cameraComponent)
}

func updateFromRegistry(r *ecs.Registry, dt float64) {
	for _, e := range ecs.View[*TimeComponent](r) {
		timeComponent, _ := ecs.Get[*TimeComponent](r, e)
		if textComponent, hasText := ecs.Get[*TextComponent](r, e); hasText {
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
