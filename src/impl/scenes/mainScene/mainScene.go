package mainScene

import (
	"image/color"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type MainUpdater struct {
	PlayerUpdater
	UIUpdater
	WorldMapUpdater
	CountryUpdater

	Window *pixelgl.Window
}

func (u *MainUpdater) Update(dt float64) {
	u.PlayerUpdater.Update(u.Window, dt)
	u.UIUpdater.Update(u.Window, dt)
	u.WorldMapUpdater.Update(u.Window, dt)
	u.CountryUpdater.Update(u.Window, dt)
}

func NewScene(win *pixelgl.Window) *ecs.Scene {
	registry := ecs.NewRegistry()
	updater := &MainUpdater{}
	mainScene := ecs.NewScene(registry, updater)
	configureScene(mainScene, updater, win)

	return mainScene
}

func NewFactory(s *ecs.Scene, frame int, position pixel.Vec, orig pixel.Vec) ecs.EntityFactory {
	prefab := CountryPrefab{
		Frame:    frame,
		Position: position,
		Orig:     orig,
	}
	return ecs.NewEntityFactory(s, prefab)
}

func configureScene(s *ecs.Scene, u *MainUpdater, win *pixelgl.Window) {
	mainCamera := &components.CameraComponent{}
	mainCamera.Init(1.0, 1.2, true)

	UICanvas := &components.CanvasComponent{}
	UICanvas.Init(win.Bounds(), color.RGBA{R: 0, G: 0, B: 0, A: 0})

	worldTime := &components.TimeComponent{}
	worldTime.Init("Mon, 02 Jan 2006 15:04:05 MST")

	worldClock := &components.TextComponent{}
	worldClock.Init(pixel.V(10, 10), text.NewAtlas(basicfont.Face7x13, text.ASCII))

	UI := s.CreateEntity()
	ecs.AddComponent(UI, UICanvas)
	ecs.AddComponent(UI, worldTime)
	ecs.AddComponentGroup[components.UIElement](UI, worldClock)

	Player := s.CreateEntity()
	ecs.AddComponent(Player, mainCamera)
	ecs.AddComponentGroup[components.Drawable](Player, UICanvas)

	// TODO: Temporary. Probably not going to generate initial countries here
	initialCountryFactory := NewFactory(s, 0, win.Bounds().Center(), pixel.ZV)
	initialCountry := initialCountryFactory.Generate()

	// Map the necessary entities onto the updater
	u.Window = win
	u.UI = UI
	u.Player = Player
	u.Countries = make([]ecs.Entity, 0, 1)
	u.Countries = append(u.Countries, initialCountry)
}
