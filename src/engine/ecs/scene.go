package ecs

type Updater interface {
	Update(float64)
}

type UpdaterFunc func(float64)

func (f UpdaterFunc) Update(dt float64) {
	f(dt)
}

type Scene struct {
	Updater
	registry *Registry
}

// NewScene - Creates a new scene filling in required initialization parameters
func NewScene(registry *Registry, updater Updater) *Scene {
	return &Scene{
		registry: registry,
		Updater:  updater,
	}
}

func (s *Scene) CreateEntity() Entity {
	e := NewEntity(s)

	return e
}

func Map[C ComponentData](s *Scene) map[EntityId]C {
	return View[C](s.registry)
}
