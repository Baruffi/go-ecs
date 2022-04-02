package main

import (
	_ "image/png"

	"example.com/v0/src/impl"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(impl.Run)
}
