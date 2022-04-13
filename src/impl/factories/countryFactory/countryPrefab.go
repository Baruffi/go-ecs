package countryFactory

import (
	"math"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/managers"
	"example.com/v0/src/impl/tools"
	"example.com/v0/src/queue"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type CountryPrefab struct {
	frame         int
	position      pixel.Vec
	orig          pixel.Vec
	timeLoc       string
	drawerManager *managers.DrawerManager
}

func (p *CountryPrefab) Update(frame int, position pixel.Vec, orig pixel.Vec, timeLoc string) {
	p.frame = frame
	p.position = position
	p.orig = orig
	p.timeLoc = timeLoc
}

func (p CountryPrefab) Configure(countryEntity ecs.Entity) {
	timeTag := ecs.Add[components.TagComponent](countryEntity)
	timeTag.Init(p.timeLoc)

	drawComponent := ecs.Add[components.DrawComponent](countryEntity)
	spritesheet, err := tools.LoadPicture("./assets/countryFactory/countries.png")
	if err != nil {
		panic(err)
	}
	frameSizeX := 256.0
	frameSizeY := 256.0
	spriteScale := 1.0
	drawComponent.Init(spritesheet, frameSizeX, frameSizeY, spriteScale)
	drawComponent.PrepareFrame(p.frame, p.position)

	textComponent := ecs.Add[components.TextComponent](countryEntity)
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textComponent.Init(p.position.Add(p.orig), atlas, colornames.Black, 2)

	frameScaleStep := math.Sqrt(frameSizeX*frameSizeY) * spriteScale / 2
	hoverComponent := ecs.Add[components.ColliderComponent](countryEntity)
	area := pixel.R(p.position.X-frameScaleStep, p.position.Y-frameScaleStep, p.position.X+frameScaleStep, p.position.Y+frameScaleStep)
	hoverComponent.Init(area, p.position, 1.0, 0.0, 1.0)

	p.drawerManager.Enqueue(queue.FIVE, countryEntity, drawComponent)
	p.drawerManager.Enqueue(queue.FOUR, countryEntity, textComponent)
	p.drawerManager.Enqueue(queue.SEVEN, countryEntity, hoverComponent)
}
