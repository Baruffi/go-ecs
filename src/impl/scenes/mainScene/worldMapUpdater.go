package mainScene

import (
	"example.com/v0/src/ecs"
	"github.com/faiface/pixel/pixelgl"
)

type WorldMapUpdater struct {
	WorldMap ecs.Entity
}

func (u *WorldMapUpdater) Update(w *pixelgl.Window, dt float64) {
}
