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

func (s *Scene) MoveEntity(e Entity) {
	e.JoinScene(s)
}

func Map[C ComponentData](s *Scene, is ...TypedComponentId[C]) map[EntityId]C {
	// Since only 1 id per type will exist in the registry, there should be no use case with multiple ids as args. Using ... as an optional notation
	for _, i := range is {
		return ViewById(s.registry, i)
	}
	return View[C](s.registry)
}

func MapGroup[C ComponentData](s *Scene, is ...TypedComponentGroupId[C]) []map[EntityId]C {
	// Since only 1 id per type will exist in the registry, there should be no use case with multiple ids as args. Using ... as an optional notation
	for _, i := range is {
		return ViewGroupById(s.registry, i)
	}
	return ViewGroup[C](s.registry)
}
