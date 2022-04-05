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

type TransformableComponent struct {
	Matrix pixel.Matrix
}

func (t *TransformableComponent) Transform(ref pixel.Vec, delta pixel.Vec, deltaZoom float64) {
	t.Matrix = t.Matrix.Chained(pixel.IM.Scaled(ref, deltaZoom).Moved(delta))
}

func (t *TransformableComponent) InverseTransform(ref pixel.Vec, delta pixel.Vec, deltaZoom float64, totalZoom float64) {
	t.Matrix = t.Matrix.Chained(pixel.IM.Moved(pixel.ZV.Sub(delta.Scaled(1/totalZoom))).Scaled(ref, 1/deltaZoom))
}

func (t *TransformableComponent) Project(position pixel.Vec) pixel.Vec {
	return t.Matrix.Project(position)
}

func (t *TransformableComponent) Unproject(position pixel.Vec) pixel.Vec {
	return t.Matrix.Unproject(position)
}

type CameraComponent struct {
	TransformableComponent
	CamDelta     pixel.Vec
	CamDeltaZoom float64
	CamZoom      float64
	CamZoomSpeed float64
	Active       bool
}

func (c *CameraComponent) Init(zoom float64, zoomSpeed float64, active bool) {
	c.TransformableComponent = TransformableComponent{
		Matrix: pixel.IM,
	}
	c.CamDelta = pixel.ZV
	c.CamDeltaZoom = 1.0
	c.CamZoom = zoom
	c.CamZoomSpeed = zoomSpeed
	c.Active = active
}

func (c *CameraComponent) Toggle() {
	c.Active = !c.Active
}

func (c *CameraComponent) Scroll(scroll pixel.Vec) {
	c.CamDeltaZoom = math.Pow(c.CamZoomSpeed, scroll.Y)
	c.CamZoom *= c.CamDeltaZoom
}

func (c *CameraComponent) Move(delta pixel.Vec) {
	c.CamDelta = delta
}

func (c *CameraComponent) Update(center pixel.Vec) {
	c.Transform(center, c.CamDelta, c.CamDeltaZoom)
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

func (c *CanvasComponent) Clear() {
	c.Canvas.Clear(c.Color)
}

func (c *CanvasComponent) Draw(target pixel.Target) {
	c.Canvas.Draw(target, c.Matrix)
}

type DrawComponent struct {
	Batch       *pixel.Batch
	Spritesheet pixel.Picture
	FrameSize   float64
	SpriteScale float64
	Frames      []pixel.Rect
}

func (d *DrawComponent) Init(spritesheet pixel.Picture, frameSize float64, spriteScale float64) {
	d.Spritesheet = spritesheet
	d.FrameSize = frameSize
	d.SpriteScale = spriteScale
	d.Frames = make([]pixel.Rect, 0, int(spritesheet.Bounds().Norm().Area()/frameSize))
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += frameSize {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += frameSize {
			d.Frames = append(d.Frames, pixel.R(x, y, x+frameSize, y+frameSize))
		}
	}

	d.Batch = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
}

func (d *DrawComponent) PrepareFrame(frame int, position pixel.Vec) error {
	if frame >= len(d.Frames) {
		return errors.New("frame out of bounds")
	}
	sprite := pixel.NewSprite(d.Spritesheet, d.Frames[frame])
	sprite.Draw(d.Batch, pixel.IM.Scaled(pixel.ZV, d.SpriteScale).Moved(position))

	return nil
}

func (d *DrawComponent) ClearFrames() {
	d.Batch.Clear()
}

func (d *DrawComponent) Draw(target pixel.Target) {
	d.Batch.Draw(target)
}

func (d *DrawComponent) DrawFrame(frame int, position pixel.Vec, target pixel.Target) error {
	d.Batch.Clear()
	err := d.PrepareFrame(frame, position)
	d.Batch.Draw(target)

	return err
}

type TextComponent struct {
	Txt *text.Text
}

func (t *TextComponent) Init(orig pixel.Vec, atlas *text.Atlas) {
	t.Txt = text.New(orig, atlas)
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

type HoverComponent struct {
	Area      pixel.Rect
	IsHovered bool
}

func (h *HoverComponent) Hover(mousePosition pixel.Vec) {
	h.IsHovered = h.Area.Contains(mousePosition)
}
