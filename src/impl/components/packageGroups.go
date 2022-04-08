package components

import "github.com/faiface/pixel"

type Drawer interface {
	Draw(pixel.Target)
}
