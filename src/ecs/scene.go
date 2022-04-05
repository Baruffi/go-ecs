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

func Map[D ComponentData](s *Scene, is ...TypedComponentId[D]) map[EntityId]D {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return ViewById(s.registry, is[0])
	default:
		return View[D](s.registry)
	}
}

func MapGroup[D ComponentData](s *Scene, is ...TypedComponentGroupId[D]) []map[EntityId]D {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return ViewGroupById(s.registry, is[0])
	default:
		return ViewGroup[D](s.registry)
	}
}

func ClearComponentType[D ComponentData](s *Scene, is ...TypedComponentId[D]) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		ClearTypeById(s.registry, is[0])
	default:
		ClearType[D](s.registry)
	}
}

func ClearComponentGroup[D ComponentData](s *Scene, is ...TypedComponentGroupId[D]) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		ClearGroupById(s.registry, is[0])
	default:
		ClearGroup[D](s.registry)
	}
}
