package impl

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

type TextComponent struct {
	txt *text.Text
}

func (t *TextComponent) Init(orig pixel.Vec, atlas *text.Atlas) {
	t.txt = text.New(orig, atlas)
}

func (t *TextComponent) Write(str string) {
	fmt.Fprintln(t.txt, str)
}

func (t *TextComponent) Clear() {
	t.txt.Clear()
}

func (t *TextComponent) Draw(surface pixel.Target) {
	t.txt.Draw(surface, pixel.IM.Scaled(t.txt.Orig, 2))
}

type DrawComponent struct {
	batch       *pixel.Batch
	spritesheet pixel.Picture
	frameSize   float64
	spriteScale float64
	frames      []pixel.Rect
}

func (d *DrawComponent) Init(spritesheet pixel.Picture, frameSize float64, spriteScale float64) {
	d.spritesheet = spritesheet
	d.frameSize = frameSize
	d.spriteScale = spriteScale
	d.frames = make([]pixel.Rect, 0, int(spritesheet.Bounds().Norm().Area()/frameSize))
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += frameSize {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += frameSize {
			d.frames = append(d.frames, pixel.R(x, y, x+frameSize, y+frameSize))
		}
	}

	d.batch = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
}

func (d *DrawComponent) PrepareFrame(frame int, position pixel.Vec) error {
	if frame >= len(d.frames) {
		return errors.New("frame out of bounds")
	}
	sprite := pixel.NewSprite(d.spritesheet, d.frames[frame])
	sprite.Draw(d.batch, pixel.IM.Scaled(pixel.ZV, d.spriteScale).Moved(position))

	return nil
}

func (d *DrawComponent) ClearFrames() {
	d.batch.Clear()
}

func (d *DrawComponent) Draw(surface pixel.Target) {
	d.batch.Draw(surface)
}

func (d *DrawComponent) DrawFrame(frame int, position pixel.Vec, surface pixel.Target) error {
	d.batch.Clear()
	err := d.PrepareFrame(frame, position)
	d.batch.Draw(surface)

	return err
}

type CameraComponent struct {
	cam          pixel.Matrix
	camPos       pixel.Vec
	camSpeed     float64
	camZoom      float64
	camZoomSpeed float64
}

func (c *CameraComponent) Update(position pixel.Vec, scroll pixel.Vec) {
	c.cam = pixel.IM.Scaled(c.camPos, c.camZoom).Moved(position.Sub(c.camPos))
	c.camZoom *= math.Pow(c.camZoomSpeed, scroll.Y)
}

func (c *CameraComponent) Unproject(position pixel.Vec) pixel.Vec {
	return c.cam.Unproject(position)
}

type TimeComponent struct {
	time   time.Time
	ticker *time.Ticker
	format string
}

func (t *TimeComponent) Format() string {
	return t.time.Local().Format(t.format)
}
