package mainScene

import (
	"image/color"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/factories/countryFactory"
	"example.com/v0/src/impl/managers"
	"example.com/v0/src/impl/tools"
	"example.com/v0/src/queue"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func NewScene(win *pixelgl.Window, eventManager *managers.EventManager, drawerManager *managers.DrawerManager) *ecs.Scene {
	mainScene := ecs.NewScene[MainUpdater]()
	configureScene(mainScene, win, eventManager, drawerManager)

	return mainScene
}

func configureScene(s *ecs.Scene, win *pixelgl.Window, eventManager *managers.EventManager, drawerManager *managers.DrawerManager) {
	UI := s.CreateEntity()
	UICanvas := ecs.Add[components.CanvasComponent](UI)
	clock := ecs.Add[components.Combiner[components.TimeComponent, components.TextComponent]](UI)

	world := s.CreateEntity()
	worldMap := ecs.Add[components.Combiner[components.DrawComponent, components.ColliderComponent]](world)

	player := s.CreateEntity()
	camera := ecs.Add[components.Combiner[components.CameraComponent, components.ColliderComponent]](player)

	cameraMatrix := components.CameraComponent{}
	cameraMatrix.Init(1.0, 1.2, true)
	cameraCollider := components.ColliderComponent{}
	cameraCollider.Init(win.Bounds(), win.Bounds().Center(), 1.0, 1.0, 1.2)
	camera.T1 = cameraMatrix
	camera.T2 = cameraCollider

	UICanvas.Init(win.Bounds(), color.RGBA{R: 0, G: 0, B: 0, A: 0})

	clockTime := components.TimeComponent{}
	clockTime.Init("UTC", "Mon, 02 Jan 2006 15:04:05 MST")
	clockText := components.TextComponent{}
	clockText.Init(pixel.V(10, 10), text.NewAtlas(basicfont.Face7x13, text.ASCII), colornames.Black, 1)
	clock.T1 = clockTime
	clock.T2 = clockText

	worldMapBackdrop := components.DrawComponent{}
	spritesheet, err := tools.LoadPicture("./assets/mainScene/A_large_blank_world_map_with_oceans_marked_in_blue.png")
	if err != nil {
		panic(err)
	}
	worldMapBackdrop.Init(spritesheet, spritesheet.Bounds().Norm().W(), spritesheet.Bounds().Norm().H(), 1)
	sprite, _ := worldMapBackdrop.PrepareFrame(0, pixel.ZV)
	worldMapCollider := components.ColliderComponent{}
	worldMapCollider.Init(sprite.Frame(), pixel.ZV, worldMapBackdrop.SpriteScale, 0.0, 1.0)
	worldMap.T1 = worldMapBackdrop
	worldMap.T2 = worldMapCollider

	// TODO: Temporary. Probably not going to generate initial countries here
	initialCountryFactory := countryFactory.NewFactory(s, 0, win.Bounds().Center(), pixel.ZV, "EST", eventManager, drawerManager)
	initialCountry := initialCountryFactory.Generate()
	initialCountryFactory.Prefab.Update(0, pixel.V(-100, -100), pixel.ZV, "MST")
	secondCountry := initialCountryFactory.Generate()
	countries := []ecs.Entity{initialCountry, secondCountry}

	// Map every component that will be always drawn
	drawerManager.Enqueue(queue.TWO, true, UICanvas)
	drawerManager.Enqueue(queue.SEVEN, true, worldMap.GetSecond(), camera.GetSecond())
	drawerManager.Enqueue(queue.NINE, true, worldMap.GetFirst())

	// Map the necessary entities onto the updater
	u := s.Updater.(MainUpdater)
	u.UI = UI
	u.World = world
	u.Countries = countries
	u.Player = player
	u.Window = win
	u.EventManager = eventManager
	u.DrawerManager = drawerManager
	s.Updater = u
}
