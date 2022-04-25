package mainScene

import (
	"fmt"
	"log"
	"math/rand"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/factories/countryFactory"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MainUpdater struct {
	pendingCountryCounter int
	countryFactoryHolder  ecs.EntityFactory[countryFactory.CountryPrefab]
	window                *pixelgl.Window

	player    ecs.Entity
	world     ecs.Entity
	ui        ecs.Entity
	countries []ecs.Entity
}

func NewUpdater(countryFactoryHolder ecs.EntityFactory[countryFactory.CountryPrefab], window *pixelgl.Window, player ecs.Entity, world ecs.Entity, ui ecs.Entity) *MainUpdater {
	updater := &MainUpdater{
		countryFactoryHolder: countryFactoryHolder,
		window:               window,

		player:    player,
		world:     world,
		ui:        ui,
		countries: make([]ecs.Entity, 0),
	}

	return updater
}

func (u *MainUpdater) GenerateCountries() {
	var timeLoc string
	switch rand.Intn(2) {
	case 1:
		timeLoc = "MST"
	case 0:
		timeLoc = "EST"
	}
	randV := func() pixel.Vec {
		return pixel.V(rand.Float64()*1000, rand.Float64()*1000)
	}
	u.countryFactoryHolder.Prefab.Update(0, randV().Sub(randV()), pixel.ZV.Sub(pixel.V(100, 0)), timeLoc)
	u.countries = append(u.countries, u.countryFactoryHolder.Generate())
}

func (u *MainUpdater) DestroyCountries() {
	for _, country := range u.countries {
		country.Die()
	}
	u.countries = make([]ecs.Entity, 0)
	u.pendingCountryCounter = 0
}

func (u *MainUpdater) Update(stage *ecs.Stage, dt float64) {
	u.pendingCountryCounter++
	if u.pendingCountryCounter%100 == 0 {
		u.GenerateCountries()
	}
	if u.pendingCountryCounter > 1000 {
		u.DestroyCountries()
	}

	// Experimental use of view
	scene := stage.GetScene()
	tagView := ecs.View(scene, &components.TagComponent{})
	for tagView.Next() {
		entity := tagView.GetEntity()
		tag := tagView.GetComponents()[0].(*components.TagComponent)
		log.Printf("tag '%s' for entity '%v'\n", tag.Tag, entity)
	}

	if clock, ok := ecs.Get[components.Combiner[components.TimeComponent, components.TextComponent]](u.ui); ok {
		timeComponent := clock.GetFirst()
		textComponent := clock.GetSecond()
		select {
		case <-timeComponent.Ticker.C:
			textComponent.Clear()
			timeComponent.UpdateTime()
			timeStr := fmt.Sprintf("TIME: %s", timeComponent.String())
			textComponent.Write(timeStr)
			if UICanvas, ok := ecs.Get[components.CanvasComponent](u.ui); ok {
				UICanvas.Clear()
				textComponent.Draw(UICanvas.Canvas)
			}
		default:
		}
	}
	if camera, ok := ecs.Get[components.Combiner[components.CameraComponent, components.ColliderComponent]](u.player); ok {
		cameraComponent := camera.GetFirst()
		cameraCollider := camera.GetSecond()
		if cameraComponent.Active {
			leftClickHeld := u.window.Pressed(pixelgl.MouseButtonLeft)
			mousePosition := u.window.MousePosition()
			mousePreviousPosition := u.window.MousePreviousPosition()
			mouseScroll := u.window.MouseScroll()

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

			if worldMap, ok := ecs.Get[components.Combiner[components.DrawComponent, components.ColliderComponent]](u.world); ok {
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

					u.window.SetMatrix(cameraComponent.Matrix)

					if UICanvas, ok := ecs.Get[components.CanvasComponent](u.ui); ok {
						UICanvas.InverseTransform(cameraComponent.Unproject(mousePosition), cameraComponent.DeltaPos, cameraComponent.DeltaScale, cameraComponent.Scale)
					}
				} else {
					cameraCollider.Area = previousArea
					cameraCollider.Scale = previousScale
					cameraCollider.Pos = previousPos
				}
			}

			for _, country := range u.countries {
				if hoverComponent, ok := ecs.Get[components.ColliderComponent](country); ok {
					if textComponent, ok := ecs.Get[components.TextComponent](country); ok {
						textComponent.Clear()
						if hoverComponent.CollidesVec(cameraComponent.Unproject(mousePosition)) {
							if timeTag, ok := ecs.Get[components.TagComponent](country); ok {
								if clock, ok := ecs.Get[components.Combiner[components.TimeComponent, components.TextComponent]](u.ui); ok {
									timeComponent := clock.GetFirst()
									timeComponent.UpdateLocation(timeTag.Tag)
								}
							}
							textComponent.Write(fmt.Sprintf("COUNTRY %v", country))
						}
					}
				}
			}
		}
	}
}
