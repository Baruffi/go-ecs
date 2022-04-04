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

type CanvasComponent struct {
	Canvas *pixelgl.Canvas
	Color  color.RGBA
	Center pixel.Vec
	Offset pixel.Vec
	Scale  float64
}

func (c *CanvasComponent) Init(bounds pixel.Rect, color color.RGBA, center pixel.Vec, camPos pixel.Vec, camZoom float64) {
	c.Canvas = pixelgl.NewCanvas(bounds)
	c.Color = color

	c.Center = center
	c.Offset = center.Sub(camPos)
	c.Scale = 1 / camZoom
}

func (c *CanvasComponent) Transform(center pixel.Vec, camPos pixel.Vec, camZoom float64) {
	c.Center = center
	c.Offset = center.Sub(camPos)
	c.Scale = 1 / camZoom
}

func (c *CanvasComponent) Clear() {
	c.Canvas.Clear(c.Color)
}

func (c *CanvasComponent) Draw(target pixel.Target) {
	c.Canvas.Draw(target, pixel.IM.Moved(c.Offset).Scaled(c.Center, c.Scale))
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

type CameraComponent struct {
	Cam          pixel.Matrix
	CamPos       pixel.Vec
	CamSpeed     float64
	CamZoom      float64
	CamZoomSpeed float64
	Active       bool
}

func (c *CameraComponent) Toggle() {
	c.Active = !c.Active
}

func (c *CameraComponent) Scroll(scroll pixel.Vec) {
	c.CamZoom *= math.Pow(c.CamZoomSpeed, scroll.Y)
}

func (c *CameraComponent) Move(delta pixel.Vec) {
	c.CamPos = c.CamPos.Add(delta)
}

func (c *CameraComponent) Update(center pixel.Vec) {
	c.Cam = pixel.IM.Scaled(center, c.CamZoom).Moved(c.CamPos)
}

func (c *CameraComponent) Project(position pixel.Vec) pixel.Vec {
	return c.Cam.Project(position)
}

func (c *CameraComponent) Unproject(position pixel.Vec) pixel.Vec {
	return c.Cam.Unproject(position)
}

type TimeComponent struct {
	Time   time.Time
	Ticker *time.Ticker
	Format string
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
