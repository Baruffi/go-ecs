package components

import "github.com/faiface/pixel"

type Drawable interface {
	Draw(pixel.Target)
}
