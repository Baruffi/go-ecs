package managers

import (
	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/components"
	"github.com/faiface/pixel/pixelgl"
)

type DrawerManager struct {
	ecs.PriorityQueueManager[components.Drawer]
}

func NewDrawerManager(window *pixelgl.Window) DrawerManager {
	return DrawerManager{
		PriorityQueueManager: ecs.NewPriorityQueueManager[components.Drawer](
			&drawerHandler{
				window: window,
			},
		),
	}
}

type drawerHandler struct {
	window *pixelgl.Window
}

func (h *drawerHandler) Handle(pkg ecs.PriorityPackage[components.Drawer]) {
	for _, drawer := range pkg.Element {
		drawer.Draw(h.window)
	}
}
