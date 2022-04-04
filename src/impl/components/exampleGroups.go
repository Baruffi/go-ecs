package components

import "github.com/faiface/pixel"

type UIElement interface {
	Draw(pixel.Target)
}

type Drawable interface {
	Draw(pixel.Target)
}
