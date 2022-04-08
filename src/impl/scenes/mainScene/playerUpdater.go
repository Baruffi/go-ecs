package mainScene

import (
	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type PlayerUpdater struct {
	Player ecs.Entity
	World  ecs.Entity
	UI     ecs.Entity
}

func (u *PlayerUpdater) Update(w *pixelgl.Window, dt float64) {
	if camera, ok := ecs.GetComponent[*components.Combiner[*components.CameraComponent, *components.ColliderComponent]](u.Player); ok {
		cameraComponent := camera.GetFirst()

		if cameraComponent.Active {
			mousePosition := w.MousePosition()
			mousePreviousPosition := w.MousePreviousPosition()
			mouseScroll := w.MouseScroll()

			if w.Pressed(pixelgl.MouseButtonLeft) {
				mouseDelta := mousePosition.Sub(mousePreviousPosition)
				cameraComponent.Move(mouseDelta)
			} else {
				cameraComponent.Move(pixel.ZV)
			}
			cameraComponent.Grow(mouseScroll.Y)
			cameraComponent.Update(mousePosition)

			w.SetMatrix(cameraComponent.Matrix)

			if UICanvas, ok := ecs.GetComponent[*components.CanvasComponent](u.UI); ok {
				UICanvas.InverseTransform(cameraComponent.Unproject(mousePosition), cameraComponent.DeltaPos, cameraComponent.DeltaScale, cameraComponent.Scale)
			}
		}
	}
}
