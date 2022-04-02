package ecs

type Surface interface {
}

type Drawer interface {
	Draw(Surface)
}

type DrawerFunc func(Surface)

func (f DrawerFunc) Draw(s Surface) {
	f(s)
}

type Renderer interface {
	CanRender() bool
	Clear()
	BeginScene()
	EndScene()
}
