package ecs

type Updater interface {
	Update(*Registry, float64)
}

type UpdaterFunc func(*Registry, float64)

func (f UpdaterFunc) Update(r *Registry, dt float64) {
	f(r, dt)
}

type Scene struct {
	registry *Registry
	updater  Updater
}

func NewScene(registry *Registry, updater Updater) *Scene {
	return &Scene{
		registry: registry,
		updater:  updater,
	}
}

func (s *Scene) CreateEntity() Entity {
	e := NewEntity(s)

	return e
}

func (s *Scene) Update(dt float64) {
	s.updater.Update(s.registry, dt)
}
