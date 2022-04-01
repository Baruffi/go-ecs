package main

import (
	"errors"
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"example.com/v0/src/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

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
	frames   int
	dt       float64
	last     time.Time
	updateCh <-chan time.Time
}

func (t *TimeComponent) Update() {
	t.dt = time.Since(t.last).Seconds()
	t.last = time.Now()
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func setup() (*ecs.Registry, []ecs.Entity) {
	registry := ecs.NewRegistry()

	timeData := TimeComponent{
		updateCh: time.Tick(time.Second),
	}
	cameraData := CameraComponent{
		camPos:       pixel.ZV,
		camSpeed:     500.0,
		camZoom:      1.0,
		camZoomSpeed: 1.2,
	}
	canvasDrawComponent := DrawComponent{}
	characterDrawComponent := DrawComponent{}

	canvas := ecs.NewEntity()

	ecs.Link(registry, canvas, &timeData)
	ecs.Link(registry, canvas, &canvasDrawComponent)

	character := ecs.NewEntity()

	ecs.Link(registry, character, &cameraData)
	ecs.Link(registry, character, &characterDrawComponent)

	return registry, []ecs.Entity{canvas, character}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	registry, entities := setup()
	canvas, character := entities[0], entities[1]

	timeComponent, _ := ecs.GetFrom[*TimeComponent](registry, canvas)
	cameraComponent, _ := ecs.GetFrom[*CameraComponent](registry, character)
	canvasDrawComponent, _ := ecs.GetFrom[*DrawComponent](registry, canvas)
	characterDrawComponent, _ := ecs.GetFrom[*DrawComponent](registry, character)

	timeComponent.Update()

	canvasSpreadsheet, err := loadPicture("trees.png")
	if err != nil {
		panic(err)
	}
	characterSpreadsheet, err := loadPicture("cursed_flushed.png")
	if err != nil {
		panic(err)
	}

	canvasDrawComponent.Init(canvasSpreadsheet, 32, 4)
	characterDrawComponent.Init(characterSpreadsheet, 256, 1)

	for !win.Closed() {
		timeComponent.Update()

		cameraComponent.Update(win.Bounds().Center(), win.MouseScroll())
		win.SetMatrix(cameraComponent.cam)

		// Clear window buffer
		win.Clear(colornames.Forestgreen)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			canvasDrawComponent.PrepareFrame(rand.Intn(len(canvasDrawComponent.frames)), cameraComponent.Unproject(win.MousePosition()))
		}
		canvasDrawComponent.Draw(win)
		characterDrawComponent.DrawFrame(0, cameraComponent.Unproject(win.Bounds().Center()), win)

		// Swap window buffers
		win.Update()

		timeComponent.frames++
		select {
		case <-timeComponent.updateCh:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, timeComponent.frames))
			timeComponent.frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
