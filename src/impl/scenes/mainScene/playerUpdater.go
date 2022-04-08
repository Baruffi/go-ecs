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
		cameraCollider := camera.GetSecond()
		if cameraComponent.Active {
			leftClickHeld := w.Pressed(pixelgl.MouseButtonLeft)
			mousePosition := w.MousePosition()
			mousePreviousPosition := w.MousePreviousPosition()
			mouseScroll := w.MouseScroll()

			previousArea := cameraCollider.Area
			previousScale := cameraCollider.Scale
			previousPos := cameraCollider.Pos

			if leftClickHeld {
				mouseDelta := mousePosition.Sub(mousePreviousPosition)
				cameraCollider.Move(pixel.ZV.Sub(mouseDelta.Scaled(cameraCollider.Scale)))
			} else {
				cameraCollider.Move(pixel.ZV)
			}

			cameraCollider.Grow(-mouseScroll.Y)
			cameraCollider.Update(cameraComponent.Unproject(mousePosition))

			if worldMap, ok := ecs.GetComponent[*components.Combiner[*components.DrawComponent, *components.ColliderComponent]](u.World); ok {
				worldMapCollider := worldMap.GetSecond()

				if worldMapCollider.CollidesVec(cameraCollider.Area.Min) && worldMapCollider.CollidesVec(cameraCollider.Area.Max) {
					if leftClickHeld {
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
				} else {
					cameraCollider.Area = previousArea
					cameraCollider.Scale = previousScale
					cameraCollider.Pos = previousPos
				}
			}
		}
	}
}
