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

func (u *UIUpdater) Update(w *pixelgl.Window, player ecs.Entity, dt float64) {
	for _, UIElement := range ecs.GetComponentGroup[components.UIElement](u.UI) {
		switch UIComponent := UIElement.(type) {
		case *components.TextComponent:
			if cameraComponent, ok := ecs.GetComponent[*components.CameraComponent](player); ok {
				// TODO: verify why this is abnormally slow to reflect on screen
				UIComponent.Txt.Orig = cameraComponent.Cam.Unproject(w.Bounds().Min)
			}
			if timeComponent, ok := ecs.GetComponent[*components.TimeComponent](u.UI); ok {
				select {
				case <-timeComponent.Ticker.C:
					UIComponent.Clear()
					timeComponent.Time = time.Now()
					timeStr := fmt.Sprintf("TIME: %s", timeComponent.String())
					UIComponent.Write(timeStr)
				default:
				}
			}
		}
	}
}
