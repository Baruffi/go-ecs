package systems

import (
	"example.com/v0/src/ecs"
	"example.com/v0/src/queue"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Drawer interface {
	Draw(pixel.Target)
}

type DrawerWrapper struct {
	owner  ecs.Entity
	drawer Drawer
}

type DrawSystem struct {
	drawerQueue   *queue.PriorityQueue[DrawerWrapper]
	drawerHandler *queue.GenericQueueHandler[DrawerWrapper]
}

func NewDrawSystem(window *pixelgl.Window) *DrawSystem {
	return &DrawSystem{
		drawerQueue: queue.NewPriorityQueue[DrawerWrapper](),
		drawerHandler: queue.NewGenericQueueHandler[DrawerWrapper](&drawerHandler{
			window: window,
		}),
	}
}

func (m *DrawSystem) Enqueue(level queue.PriorityLevel, owner ecs.Entity, drawers ...Drawer) error {
	var err error
	m.drawerQueue.SetEnqueueLevel(level)
	for _, drawer := range drawers {
		err = m.drawerQueue.Enqueue(DrawerWrapper{owner, drawer})
	}
	return err
}

func (m *DrawSystem) Execute() error {
	return m.drawerHandler.Consume(m.drawerQueue)
}

type drawerHandler struct {
	window *pixelgl.Window
}

func (h *drawerHandler) Handle(wrapper DrawerWrapper) queue.HandlerResult {
	wrapper.drawer.Draw(h.window)

	if wrapper.owner.IsAlive() {
		return queue.NOT_DONE
	}
	return queue.DONE
}
