package mainScene

import (
	"math"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/managers"
	"example.com/v0/src/impl/scenes"
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
	timeTag := &components.TagComponent{}
	timeTag.Init(p.timeLoc)

	drawComponent := &components.DrawComponent{}
	spritesheet, err := scenes.LoadPicture("../assets/countries.png")
	if err != nil {
		panic(err)
	}
	frameSizeX := 256.0
	frameSizeY := 256.0
	spriteScale := 1.0
	drawComponent.Init(spritesheet, frameSizeX, frameSizeY, spriteScale)
	drawComponent.PrepareFrame(p.frame, p.position)

	textComponent := &components.TextComponent{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textComponent.Init(p.position.Add(p.orig), atlas, colornames.Black, 2)

	frameScaleStep := math.Sqrt(frameSizeX*frameSizeY) * spriteScale / 2
	hoverComponent := &components.ColliderComponent{}
	area := pixel.R(p.position.X-frameScaleStep, p.position.Y-frameScaleStep, p.position.X+frameScaleStep, p.position.Y+frameScaleStep)
	hoverComponent.Init(area, p.position, 1.0, 0.0, 1.0)

	ecs.AddComponent(countryEntity, timeTag)
	ecs.AddComponent(countryEntity, drawComponent)
	ecs.AddComponent(countryEntity, textComponent)
	ecs.AddComponent(countryEntity, hoverComponent)

	p.drawerManager.Enqueue(managers.FIVE, true, drawComponent)
	p.drawerManager.Enqueue(managers.FOUR, true, textComponent)
	p.drawerManager.Enqueue(managers.SEVEN, true, hoverComponent)
}
