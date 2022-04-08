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
	*Registry
}

// NewScene - Creates a new scene filling in required initialization parameters
func NewScene(updater Updater) *Scene {
	return &Scene{
		Registry: NewRegistry(),
		Updater:  updater,
	}
}

func (s *Scene) CreateEntity() Entity {
	e := NewEntity(s)

	return e
}

func (s *Scene) ClearEntity(e Entity) {
	ClearEntity(s.Registry, e.id)
}

func (s *Scene) Clear() {
	s.Registry.Clear()
}

func Map[D ComponentData](s *Scene) map[EntityId]D {
	return View[D](s.Registry)
}

func MapGroup[D ComponentData](s *Scene) map[ComponentId]map[EntityId]D {
	return ViewGroup[D](s.Registry)
}

func ClearComponentType[D ComponentData](s *Scene) {
	ClearType[D](s.Registry)
}

func ClearComponentGroup[D ComponentData](s *Scene) {
	ClearGroup[D](s.Registry)
}
