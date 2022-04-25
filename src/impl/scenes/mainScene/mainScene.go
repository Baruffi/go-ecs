package mainScene

import (
	"image/color"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/systems"
	"example.com/v0/src/impl/tools"
	"example.com/v0/src/queue"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func NewScene() *ecs.Scene {
	mainScene := ecs.NewScene()

	return mainScene
}

func ConfigureScene(s *ecs.Scene, win *pixelgl.Window, eventSystem *systems.EventSystem, drawSystem *systems.DrawSystem) (player ecs.Entity, world ecs.Entity, ui ecs.Entity) {
	ui = s.CreateEntity()
	UICanvas := ecs.Add[components.CanvasComponent](ui)
	clock := ecs.Add[components.Combiner[components.TimeComponent, components.TextComponent]](ui)

	world = s.CreateEntity()
	worldMap := ecs.Add[components.Combiner[components.DrawComponent, components.ColliderComponent]](world)

	player = s.CreateEntity()
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

	// Map every component that will be always drawn
	drawSystem.Enqueue(queue.SEVEN, ui, UICanvas)
	drawSystem.Enqueue(queue.TWO, player, worldMap.GetSecond(), camera.GetSecond())
	drawSystem.Enqueue(queue.ZERO, world, worldMap.GetFirst())

	return player, world, ui
}
