package ecs

type Type interface {
}

type Factory[T Type] interface {
	Generate() T
}

type Prefab[T Type] interface {
	Configure(T)
}

type EntityFactory struct {
	Prefab[Entity]
	s *Scene
}

func NewEntityFactory(s *Scene, prefab Prefab[Entity]) EntityFactory {
	return EntityFactory{
		s:      s,
		Prefab: prefab,
	}
}

func (f *EntityFactory) Generate() Entity {
	e := f.s.CreateEntity()
	f.Configure(e)
	return e
}
