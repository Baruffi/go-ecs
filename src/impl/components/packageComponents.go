package components

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type CameraComponent struct {
	DeltaComponent
	TransformableComponent
	ActiveComponent
}

func (c *CameraComponent) Init(initialZoom float64, zoomSpeed float64, active bool) {
	c.TransformableComponent.Init()
	c.DeltaComponent.Init(pixel.ZV, initialZoom, 1, zoomSpeed)
	c.ActiveComponent.Init(active)
}

func (c *CameraComponent) Update(center pixel.Vec) {
	c.DeltaComponent.Update()
	c.Transform(center, c.DeltaPos, c.DeltaScale)
}

type CanvasComponent struct {
	TransformableComponent
	Canvas *pixelgl.Canvas
	Color  color.RGBA
}

func (c *CanvasComponent) Init(bounds pixel.Rect, color color.RGBA) {
	c.Canvas = pixelgl.NewCanvas(bounds)
	c.Matrix = pixel.IM.Moved(bounds.Center())
	c.Color = color
}

func (c *CanvasComponent) GetLayer() Layer {
	return Layer2
}

func (c *CanvasComponent) Bounds() pixel.Rect {
	return c.Canvas.Bounds()
}

func (c *CanvasComponent) Clear() {
	c.Canvas.Clear(c.Color)
}

func (c *CanvasComponent) Draw(target pixel.Target) {
	c.Canvas.Draw(target, c.Matrix)
}

type DrawComponent struct {
	Batch       *pixel.Batch
	Spritesheet pixel.Picture
	SpriteScale float64
	Frames      []pixel.Rect
	Layer       Layer
}

func (d *DrawComponent) Init(layer Layer, spritesheet pixel.Picture, frameSizeX float64, frameSizeY float64, spriteScale float64) {
	d.Layer = layer
	d.Spritesheet = spritesheet
	d.SpriteScale = spriteScale
	frameArea := math.Sqrt(frameSizeX * frameSizeY)
	d.Frames = make([]pixel.Rect, 0, int(spritesheet.Bounds().Norm().Area()/frameArea))
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += frameSizeX {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += frameSizeY {
			d.Frames = append(d.Frames, pixel.R(x, y, x+frameSizeX, y+frameSizeY))
		}
	}
	d.Batch = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
}

func (d *DrawComponent) GetLayer() Layer {
	return d.Layer
}

func (d *DrawComponent) PrepareFrame(frame int, position pixel.Vec) (*pixel.Sprite, error) {
	if frame >= len(d.Frames) {
		return nil, errors.New("frame out of bounds")
	}
	sprite := pixel.NewSprite(d.Spritesheet, d.Frames[frame])
	sprite.Draw(d.Batch, pixel.IM.Scaled(sprite.Frame().Center(), d.SpriteScale).Moved(position))

	return sprite, nil
}

func (d *DrawComponent) Clear() {
	d.Batch.Clear()
}

func (d *DrawComponent) Draw(target pixel.Target) {
	d.Batch.Draw(target)
}

func (d *DrawComponent) DrawFrame(frame int, position pixel.Vec, target pixel.Target) (*pixel.Sprite, error) {
	d.Clear()
	sprite, err := d.PrepareFrame(frame, position)
	d.Draw(target)

	return sprite, err
}

type TextComponent struct {
	Txt *text.Text
}

func (t *TextComponent) Init(orig pixel.Vec, atlas *text.Atlas, color color.RGBA) {
	t.Txt = text.New(orig, atlas)
	t.Txt.Color = color
}

func (t *TextComponent) GetLayer() Layer {
	return Layer0
}

func (t *TextComponent) Write(str string) {
	fmt.Fprintln(t.Txt, str)
}

func (t *TextComponent) Clear() {
	t.Txt.Clear()
}

func (t *TextComponent) Draw(target pixel.Target) {
	t.Txt.Draw(target, pixel.IM)
}

type TimeComponent struct {
	Time   time.Time
	Ticker *time.Ticker
	Format string
}

func (t *TimeComponent) Init(format string) {
	t.Ticker = time.NewTicker(time.Second)
	t.Format = format
}

func (t *TimeComponent) String() string {
	return t.Time.Local().Format(t.Format)
}
