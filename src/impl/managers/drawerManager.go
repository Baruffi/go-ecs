package managers

import (
	"example.com/v0/src/queue"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Drawer interface {
	Draw(pixel.Target)
}

type DrawerWrapper struct {
	isPersistent bool
	drawer       Drawer
}

type DrawerManager struct {
	drawerQueue   *queue.ThreadSafeQueue[DrawerWrapper, *queue.PriorityQueue[DrawerWrapper]]
	drawerHandler *queue.GenericQueueHandler[DrawerWrapper]
}

func NewDrawerManager(window *pixelgl.Window) *DrawerManager {
	return &DrawerManager{
		drawerQueue: queue.NewThreadSafeQueue[DrawerWrapper](queue.NewPriorityQueue[DrawerWrapper]()),
		drawerHandler: queue.NewGenericQueueHandler[DrawerWrapper](&drawerHandler{
			window: window,
		}),
	}
}

func (m *DrawerManager) Enqueue(level queue.PriorityLevel, isPersistent bool, drawers ...Drawer) error {
	m.drawerQueue.SafeWrite(func(queue *queue.PriorityQueue[DrawerWrapper]) {
		queue.SetEnqueueLevel(level)
	})
	for _, drawer := range drawers {
		err := m.drawerQueue.Enqueue(DrawerWrapper{isPersistent, drawer})
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *DrawerManager) Execute() (handlerErr error) {
	return m.drawerHandler.Consume(m.drawerQueue)
}

type drawerHandler struct {
	window *pixelgl.Window
}

func (h *drawerHandler) Handle(wrapper DrawerWrapper) queue.HandlerResult {
	wrapper.drawer.Draw(h.window)

	if wrapper.isPersistent {
		return queue.NOT_DONE
	}
	return queue.DONE
}
