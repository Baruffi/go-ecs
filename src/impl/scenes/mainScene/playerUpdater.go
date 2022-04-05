package mainScene

import (
	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type PlayerUpdater struct {
	Player ecs.Entity
}

func (u *PlayerUpdater) Update(w *pixelgl.Window, dt float64) {
	if cameraComponent, ok := ecs.GetComponent[*components.CameraComponent](u.Player); ok {
		if cameraComponent.Active {
			if w.Pressed(pixelgl.MouseButtonLeft) {
				mouseDelta := w.MousePosition().Sub(w.MousePreviousPosition())
				cameraComponent.Move(mouseDelta)
			} else {
				cameraComponent.Move(pixel.ZV)
			}

			cameraComponent.Scroll(w.MouseScroll())
			cameraComponent.Update(w.MousePosition())

			w.SetMatrix(cameraComponent.Matrix)

			for _, drawable := range ecs.GetComponentGroup[components.Drawable](u.Player) {
				if canvasComponent, ok := drawable.(*components.CanvasComponent); ok {
					canvasComponent.InverseTransform(cameraComponent.Unproject(w.MousePosition()), cameraComponent.CamDelta, cameraComponent.CamDeltaZoom, cameraComponent.CamZoom)
				}
			}
		}
	}
}
