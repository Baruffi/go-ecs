package components

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type TagComponent struct {
	Tag string
}

func (t *TagComponent) Init(tag string) {
	t.Tag = tag
}

type TransformableComponent struct {
	Matrix pixel.Matrix
}

func (t *TransformableComponent) Init() {
	t.Matrix = pixel.IM
}

func (t *TransformableComponent) Transform(ref pixel.Vec, delta pixel.Vec, deltaScale float64) {
	t.Matrix = t.Matrix.Chained(pixel.IM.Scaled(ref, deltaScale).Moved(delta))
}

func (t *TransformableComponent) InverseTransform(ref pixel.Vec, delta pixel.Vec, deltaScale float64, totalScale float64) {
	t.Matrix = t.Matrix.Chained(pixel.IM.Moved(pixel.ZV.Sub(delta.Scaled(1/totalScale))).Scaled(ref, 1/deltaScale))
}

func (t *TransformableComponent) Project(position pixel.Vec) pixel.Vec {
	return t.Matrix.Project(position)
}

func (t *TransformableComponent) Unproject(position pixel.Vec) pixel.Vec {
	return t.Matrix.Unproject(position)
}

type DeltaComponent struct {
	DeltaPos   pixel.Vec
	Pos        pixel.Vec
	DeltaScale float64
	Scale      float64
	Speed      float64
	ScaleSpeed float64
}

func (d *DeltaComponent) Init(initialPos pixel.Vec, initialScale float64, speed float64, scaleSpeed float64) {
	d.DeltaPos = pixel.ZV
	d.Pos = initialPos
	d.DeltaScale = 1.0
	d.Scale = initialScale
	d.Speed = speed
	d.ScaleSpeed = scaleSpeed
}

func (d *DeltaComponent) Grow(delta float64) {
	d.DeltaScale = math.Pow(d.ScaleSpeed, delta)
}

func (d *DeltaComponent) Move(delta pixel.Vec) {
	d.DeltaPos = delta.Scaled(d.Speed)
}

func (d *DeltaComponent) Update() {
	d.Scale *= d.DeltaScale
	d.Pos = d.Pos.Add(d.DeltaPos.Scaled(d.Scale))
}

type ColliderComponent struct {
	DeltaComponent
	Area     pixel.Rect
	DebugImd *imdraw.IMDraw
}

func (c *ColliderComponent) Init(area pixel.Rect, initialPos pixel.Vec, initialScale float64, speed float64, scaleSpeed float64) {
	c.Area = area.Moved(initialPos.Sub(area.Center()))
	c.DeltaComponent.Init(initialPos, initialScale, speed, scaleSpeed)
	c.DebugImd = imdraw.New(nil)
}

func (c *ColliderComponent) Update(anchor pixel.Vec) {
	c.DeltaComponent.Update()
	size := c.Area.Size().Scaled(c.DeltaScale)
	c.Area = c.Area.Resized(anchor, size).Moved(c.DeltaPos)
}

func (c *ColliderComponent) CollidesVec(position pixel.Vec) bool {
	return c.Area.Contains(position)
}

func (c *ColliderComponent) CollidesRect(bounds pixel.Rect) bool {
	return c.Area.Intersects(bounds)
}

func (c *ColliderComponent) Draw(t pixel.Target) {
	c.DebugImd.Clear()

	c.DebugImd.Color = color.RGBA{255, 0, 0, 255}
	c.DebugImd.Push(c.Area.Min)
	c.DebugImd.Push(pixel.V(c.Area.Max.X, c.Area.Min.Y))
	c.DebugImd.Push(c.Area.Max)
	c.DebugImd.Push(pixel.V(c.Area.Min.X, c.Area.Max.Y))
	c.DebugImd.Polygon(5.0)

	c.DebugImd.Draw(t)
}

type ActiveComponent struct {
	Active bool
}

func (a *ActiveComponent) Init(active bool) {
	a.Active = active
}

func (a *ActiveComponent) Toggle() {
	a.Active = !a.Active
}

type DataComponent struct {
	Data any
}
