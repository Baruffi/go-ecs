package mainScene

import (
	"fmt"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/managers"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MainUpdater struct {
	Player    ecs.Entity
	Countries []ecs.Entity
	World     ecs.Entity
	UI        ecs.Entity

	Window        *pixelgl.Window
	EventManager  *managers.EventManager
	DrawerManager *managers.DrawerManager
}

func (u MainUpdater) Update(dt float64) {
	if clock, ok := ecs.Get[components.Combiner[components.TimeComponent, components.TextComponent]](u.UI); ok {
		timeComponent := clock.GetFirst()
		textComponent := clock.GetSecond()
		select {
		case <-timeComponent.Ticker.C:
			textComponent.Clear()
			timeComponent.UpdateTime()
			timeStr := fmt.Sprintf("TIME: %s", timeComponent.String())
			textComponent.Write(timeStr)
			if UICanvas, ok := ecs.Get[components.CanvasComponent](u.UI); ok {
				UICanvas.Clear()
				textComponent.Draw(UICanvas.Canvas)
			}
		default:
		}
	}
	if camera, ok := ecs.Get[components.Combiner[components.CameraComponent, components.ColliderComponent]](u.Player); ok {
		cameraComponent := camera.GetFirst()
		cameraCollider := camera.GetSecond()
		if cameraComponent.Active {
			leftClickHeld := u.Window.Pressed(pixelgl.MouseButtonLeft)
			mousePosition := u.Window.MousePosition()
			mousePreviousPosition := u.Window.MousePreviousPosition()
			mouseScroll := u.Window.MouseScroll()

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

			if worldMap, ok := ecs.Get[components.Combiner[components.DrawComponent, components.ColliderComponent]](u.World); ok {
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

					u.Window.SetMatrix(cameraComponent.Matrix)

					if UICanvas, ok := ecs.Get[components.CanvasComponent](u.UI); ok {
						UICanvas.InverseTransform(cameraComponent.Unproject(mousePosition), cameraComponent.DeltaPos, cameraComponent.DeltaScale, cameraComponent.Scale)
					}
				} else {
					cameraCollider.Area = previousArea
					cameraCollider.Scale = previousScale
					cameraCollider.Pos = previousPos
				}
			}

			for _, country := range u.Countries {
				if hoverComponent, ok := ecs.Get[components.ColliderComponent](country); ok {
					if textComponent, ok := ecs.Get[components.TextComponent](country); ok {
						if hoverComponent.CollidesVec(cameraComponent.Unproject(mousePosition)) {
							if timeTag, ok := ecs.Get[components.TagComponent](country); ok {
								if clock, ok := ecs.Get[components.Combiner[components.TimeComponent, components.TextComponent]](u.UI); ok {
									timeComponent := clock.GetFirst()
									timeComponent.UpdateLocation(timeTag.Tag)
								}
							}
							textComponent.Clear()
							textComponent.Write("TEST")
						} else {
							textComponent.Clear()
						}
					}
				}
			}
		}
	}
}
