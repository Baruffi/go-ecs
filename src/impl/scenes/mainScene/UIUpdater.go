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
	if clock, ok := ecs.GetComponent[*components.Combiner[*components.TimeComponent, *components.TextComponent]](u.UI); ok {
		u.UpdateClock(clock.GetFirst(), clock.GetSecond())
	}
}

func (u *UIUpdater) UpdateClock(timeComponent *components.TimeComponent, textComponent *components.TextComponent) {
	select {
	case <-timeComponent.Ticker.C:
		textComponent.Clear()
		timeComponent.Time = time.Now()
		timeStr := fmt.Sprintf("TIME: %s", timeComponent.String())
		textComponent.Write(timeStr)
		if UICanvas, ok := ecs.GetComponent[*components.CanvasComponent](u.UI); ok {
			UICanvas.Clear()
			textComponent.Draw(UICanvas.Canvas)
		}
	default:
	}
}
