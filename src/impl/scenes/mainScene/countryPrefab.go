package mainScene

import (
	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"example.com/v0/src/impl/scenes"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
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
	spritesheet, err := scenes.LoadPicture("src/assets/countries.png")
	if err != nil {
		panic(err)
	}
	frameSize := 256.0
	spriteScale := 1.0
	drawComponent.Init(spritesheet, frameSize, spriteScale)
	drawComponent.PrepareFrame(p.Frame, p.Position)
	textComponent := &components.TextComponent{}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textComponent.Init(p.Orig, atlas)
	frameScaleStep := frameSize * spriteScale / 2
	hoverComponent := &components.HoverComponent{
		Area: pixel.R(p.Position.X-frameScaleStep, p.Position.Y-frameScaleStep, p.Position.X+frameScaleStep, p.Position.Y+frameScaleStep),
	}

	ecs.AddComponent(countryEntity, drawComponent)
	ecs.AddComponent(countryEntity, textComponent)
	ecs.AddComponent(countryEntity, hoverComponent)
	ecs.AddComponentGroup[components.Drawable](countryEntity, drawComponent)
	ecs.AddComponentGroup[components.Drawable](countryEntity, textComponent)
}
