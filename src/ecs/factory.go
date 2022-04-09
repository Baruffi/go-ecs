package ecs

type Type interface{}

type Factory[T Type] interface {
	Generate() T
}

type Prefab[T Type] interface {
	Configure(T)
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
