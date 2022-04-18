package ecs

type Factory[T any] interface {
	Generate() T
}

type FactoryFunc[T any] func() T

func (f FactoryFunc[T]) Generate() T {
	return f()
}

type Prefab[T any] interface {
	Configure(T)
}

type PrefabFunc[T any] func(T)

func (f PrefabFunc[T]) Configure(value T) {
	f(value)
}

type EntityFactory[P Prefab[Entity]] struct {
	Prefab P
	s      *Scene
}

func NewEntityFactory[P Prefab[Entity]](s *Scene, prefab P) EntityFactory[P] {
	return EntityFactory[P]{
		s:      s,
		Prefab: prefab,
	}
}

func (f EntityFactory[P]) Generate() Entity {
	e := f.s.CreateEntity()
	f.Prefab.Configure(e)
	return e
}
