package mainScene

import (
	"example.com/v0/src/ecs"
	"github.com/faiface/pixel/pixelgl"
)

type WorldUpdater struct {
	World ecs.Entity
}

func (u *WorldUpdater) Update(w *pixelgl.Window, dt float64) {
}
