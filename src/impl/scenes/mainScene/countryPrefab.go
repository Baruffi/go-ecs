package mainScene

import (
	"math"

	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/scenes"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type CountryPrefab struct {
	Frame    int
	Position pixel.Vec
	Orig     pixel.Vec
}

func (p *CountryPrefab) Update(frame int, position pixel.Vec, orig pixel.Vec) {
	p.Frame = frame
	p.Position = position
	p.Orig = orig
}

func (p CountryPrefab) Configure(countryEntity ecs.Entity) {
	drawComponent := &components.DrawComponent{}
	spritesheet, err := scenes.LoadPicture("../assets/countries.png")
	if err != nil {
		panic(err)
	}
	frameSizeX := 256.0
	frameSizeY := 256.0
	spriteScale := 1.0
	drawComponent.Init(components.Layer4, spritesheet, frameSizeX, frameSizeY, spriteScale)
	drawComponent.PrepareFrame(p.Frame, p.Position)
	textComponent := &components.TextComponent{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textComponent.Init(p.Orig, atlas, colornames.Black)
	frameScaleStep := math.Sqrt(frameSizeX*frameSizeY) * spriteScale / 2
	hoverComponent := &components.ColliderComponent{}
	area := pixel.R(p.Position.X-frameScaleStep, p.Position.Y-frameScaleStep, p.Position.X+frameScaleStep, p.Position.Y+frameScaleStep)
	hoverComponent.Init(area, p.Position, 1.0, 0.0, 1.0)

	ecs.AddComponent(countryEntity, drawComponent)
	ecs.AddComponent(countryEntity, textComponent)
	ecs.AddComponent(countryEntity, hoverComponent)
	ecs.AddComponentGroup[components.Drawer](countryEntity, drawComponent)
	ecs.AddComponentGroup[components.Drawer](countryEntity, textComponent)
}
