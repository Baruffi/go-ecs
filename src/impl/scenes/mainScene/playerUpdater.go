package mainScene

import (
	"fmt"

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
		cameraCollider := camera.GetSecond()

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

			if w.Pressed(pixelgl.MouseButtonLeft) {
				mouseDelta := mousePosition.Sub(mousePreviousPosition)
				cameraCollider.Move(pixel.ZV.Sub(mouseDelta.Scaled(cameraCollider.Scale)))
			} else {
				cameraCollider.Move(pixel.ZV)
			}

			cameraCollider.Grow(-mouseScroll.Y)
			cameraCollider.Update(cameraComponent.Unproject(mousePosition))

			if UICanvas, ok := ecs.GetComponent[*components.CanvasComponent](u.UI); ok {
				UICanvas.InverseTransform(cameraComponent.Unproject(mousePosition), cameraComponent.DeltaPos, cameraComponent.DeltaScale, cameraComponent.Scale)
			}
		}
		if worldMap, ok := ecs.GetComponent[*components.Combiner[*components.DrawComponent, *components.ColliderComponent]](u.World); ok {
			worldMapCollider := worldMap.GetSecond()

			if !worldMapCollider.CollidesRect(cameraCollider.Area) {
				fmt.Printf("OUT OF BOUNDS\n")
			}
		}
	}
}
