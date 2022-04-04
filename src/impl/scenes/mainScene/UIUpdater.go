package mainScene

import (
	"fmt"
	"time"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"github.com/faiface/pixel/pixelgl"
)

type UIUpdater struct {
	UI ecs.Entity
}

func (u *UIUpdater) Update(w *pixelgl.Window, dt float64) {
	for _, UIElement := range ecs.GetComponentGroup[components.UIElement](u.UI) {
		// TODO: Add component names to switch between components of the same type
		switch UIComponent := UIElement.(type) {
		case *components.TextComponent:
			u.UpdateClockUIComponent(UIComponent)
		}
	}
}

func (u *UIUpdater) UpdateClockUIComponent(textComponent *components.TextComponent) {
	if worldTime, ok := ecs.GetComponent[*components.TimeComponent](u.UI); ok {
		select {
		case <-worldTime.Ticker.C:
			textComponent.Clear()
			worldTime.Time = time.Now()
			timeStr := fmt.Sprintf("TIME: %s", worldTime.String())
			textComponent.Write(timeStr)
			if UICanvas, ok := ecs.GetComponent[*components.CanvasComponent](u.UI); ok {
				UICanvas.Clear()
				textComponent.Draw(UICanvas.Canvas)
			}
		default:
		}
	}
}
