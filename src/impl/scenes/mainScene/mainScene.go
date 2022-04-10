package mainScene

import (
	"fmt"
	"image/color"
	"time"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/managers"
	"example.com/v0/src/impl/scenes"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	DebugStartDrawing managers.EventCall = iota
	DebugStopDrawing
)

type MainUpdater struct {
	PlayerUpdater
	UIUpdater
	WorldUpdater
	CountryUpdater

	*pixelgl.Window
	managers.EventManager
	managers.DrawerManager
}

func (u *MainUpdater) Update(dt float64) {
	u.PlayerUpdater.Update(u.Window, dt)
	u.UIUpdater.Update(u.Window, dt)
	u.WorldUpdater.Update(u.Window, dt)
	u.CountryUpdater.Update(u.Window, dt)

	if !u.EventManager.Executing() {
		u.EventManager.AddT2(func() {
			fmt.Printf("TEST START ON %d\n", u.EventManager.GetTaskCount())
			u.EventManager.AddT1(DebugStartDrawing)
			time.Sleep(time.Second)
			u.EventManager.AddT1(DebugStopDrawing)
			time.Sleep(time.Second)
			fmt.Printf("TEST COMPLETE ON %d\n", u.EventManager.GetTaskCount())
		})
	}
}

func NewScene(win *pixelgl.Window, eventManager managers.EventManager, drawerManager managers.DrawerManager) *ecs.Scene {
	updater := &MainUpdater{}
	mainScene := ecs.NewScene(updater)
	configureScene(mainScene, updater, win, eventManager, drawerManager)

	return mainScene
}

func NewFactory(s *ecs.Scene, frame int, position pixel.Vec, orig pixel.Vec, timeLoc string, eventManager managers.EventManager, drawerManager managers.DrawerManager) ecs.EntityFactory[CountryPrefab] {
	prefab := CountryPrefab{
		frame:         frame,
		position:      position,
		orig:          orig,
		timeLoc:       timeLoc,
		drawerManager: drawerManager,
	}
	return ecs.NewEntityFactory(s, prefab)
}

func configureScene(s *ecs.Scene, u *MainUpdater, win *pixelgl.Window, eventManager managers.EventManager, drawerManager managers.DrawerManager) {
	cameraMatrix := &components.CameraComponent{}
	cameraMatrix.Init(1.0, 1.2, true)
	cameraCollider := &components.ColliderComponent{}
	cameraCollider.Init(win.Bounds(), win.Bounds().Center(), 1.0, 1.0, 1.2)
	camera := &components.Combiner[*components.CameraComponent, *components.ColliderComponent]{
		T1: cameraMatrix,
		T2: cameraCollider,
	}

	UICanvas := &components.CanvasComponent{}
	UICanvas.Init(win.Bounds(), color.RGBA{R: 0, G: 0, B: 0, A: 0})

	clockTime := &components.TimeComponent{}
	clockTime.Init("UTC", "Mon, 02 Jan 2006 15:04:05 MST")
	clockText := &components.TextComponent{}
	clockText.Init(pixel.V(10, 10), text.NewAtlas(basicfont.Face7x13, text.ASCII), colornames.Black, 1)
	clock := &components.Combiner[*components.TimeComponent, *components.TextComponent]{
		T1: clockTime,
		T2: clockText,
	}

	worldMapBackdrop := &components.DrawComponent{}
	spritesheet, err := scenes.LoadPicture("../assets/A_large_blank_world_map_with_oceans_marked_in_blue.png")
	if err != nil {
		panic(err)
	}
	worldMapBackdrop.Init(spritesheet, spritesheet.Bounds().Norm().W(), spritesheet.Bounds().Norm().H(), 1)
	sprite, _ := worldMapBackdrop.PrepareFrame(0, pixel.ZV)
	worldMapCollider := &components.ColliderComponent{}
	worldMapCollider.Init(sprite.Frame(), pixel.ZV, worldMapBackdrop.SpriteScale, 0.0, 1.0)
	worldMap := &components.Combiner[*components.DrawComponent, *components.ColliderComponent]{
		T1: worldMapBackdrop,
		T2: worldMapCollider,
	}

	UI := s.CreateEntity()
	ecs.AddComponent(UI, UICanvas)
	ecs.AddComponent(UI, clock)

	world := s.CreateEntity()
	ecs.AddComponent(world, worldMap)

	player := s.CreateEntity()
	ecs.AddComponent(player, camera)

	// TODO: Temporary. Probably not going to generate initial countries here
	initialCountryFactory := NewFactory(s, 0, win.Bounds().Center(), pixel.ZV, "EST", eventManager, drawerManager)
	initialCountry := initialCountryFactory.Generate()
	initialCountryFactory.Prefab.Update(0, pixel.V(-100, -100), pixel.ZV, "MST")
	secondCountry := initialCountryFactory.Generate()
	countries := []ecs.Entity{initialCountry, secondCountry}

	// Map every component that will be always drawn
	drawerManager.AddDefault(ecs.Level2, UICanvas)
	drawerManager.AddDefault(ecs.Level7, worldMapCollider, cameraCollider)
	drawerManager.AddDefault(ecs.Level9, worldMapBackdrop)

	// More debug stuff
	debugImd := imdraw.New(nil)
	debugImd.Color = color.RGBA{255, 0, 0, 255}
	debugImd.Push(pixel.ZV)
	debugImd.Push(pixel.V(0, 10))
	debugImd.Push(pixel.V(10, 10))
	debugImd.Push(pixel.V(10, 0))
	debugImd.Polygon(5.0)
	eventManager.SetDefault(DebugStartDrawing, func() {
		drawerManager.AddDefault(ecs.Level0, debugImd)
	})
	eventManager.SetDefault(DebugStopDrawing, func() {
		drawerManager.UnsetDefault(ecs.Level0)
	})

	// Map the necessary entities onto the updater
	u.UIUpdater.UI = UI
	u.WorldUpdater.World = world
	u.CountryUpdater.Countries = countries
	u.PlayerUpdater.Player = player
	u.PlayerUpdater.UI = UI
	u.PlayerUpdater.World = world
	u.PlayerUpdater.Countries = countries

	u.Window = win
	u.EventManager = eventManager
	u.DrawerManager = drawerManager
}
