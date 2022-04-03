package impl

import (
	"github.com/faiface/pixel"
)

type Drawable2 interface {
	Draw(pixel.Target)
}

type Drawer struct {
	target    pixel.Target
	drawables []Drawable2
}

func (d *Drawer) Draw() {
	for _, drawable := range d.drawables {
		drawable.Draw(d.target)
	}
}
