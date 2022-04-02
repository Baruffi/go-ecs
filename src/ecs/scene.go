package ecs

import "github.com/faiface/pixel"

type Drawer interface {
	Draw(*Registry, pixel.Target)
}

type DrawerFunc func(*Registry, pixel.Target)

func (f DrawerFunc) Draw(r *Registry, t pixel.Target) {
	f(r, t)
}

type Updater interface {
	Update(*Registry, float64)
}

type UpdaterFunc func(*Registry, float64)

func (f UpdaterFunc) Update(r *Registry, dt float64) {
	f(r, dt)
}

type Scene struct {
	*Registry
	Updater
	Drawer
}

func (s *Scene) Update(dt float64) {
	s.Updater.Update(s.Registry, dt)
}

func (s *Scene) Draw(t pixel.Target) {
	s.Drawer.Draw(s.Registry, t)
}
