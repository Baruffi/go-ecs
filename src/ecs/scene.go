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

func (s *Scene) Clear() {
	s.registry.Clear()
}

func Map[C ComponentData](s *Scene, is ...TypedComponentId[C]) map[EntityId]C {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return ViewById(s.registry, is[0])
	default:
		return View[C](s.registry)
	}
}

func MapGroup[C ComponentData](s *Scene, is ...TypedComponentGroupId[C]) []map[EntityId]C {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return ViewGroupById(s.registry, is[0])
	default:
		return ViewGroup[C](s.registry)
	}
}

func ClearComponentType[C ComponentData](s *Scene, is ...TypedComponentId[C]) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		ClearTypeById(s.registry, is[0])
	default:
		ClearType[C](s.registry)
	}
}

func ClearComponentGroup[C ComponentData](s *Scene, is ...TypedComponentGroupId[C]) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		ClearGroupById(s.registry, is[0])
	default:
		ClearGroup[C](s.registry)
	}
}
