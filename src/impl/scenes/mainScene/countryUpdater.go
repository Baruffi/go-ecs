package mainScene

import (
	"example.com/v0/src/ecs"
	"github.com/faiface/pixel/pixelgl"
)

type CountryUpdater struct {
	Countries []ecs.Entity
}

func (u *CountryUpdater) Update(w *pixelgl.Window, dt float64) {
}
