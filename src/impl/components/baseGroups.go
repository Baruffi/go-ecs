package components

type Combiner[T1 any, T2 any] struct {
	T1 T1
	T2 T2
}

func (c *Combiner[T1, T2]) GetFirst() T1 {
	return c.T1
}

func (c *Combiner[T1, T2]) GetSecond() T2 {
	return c.T2
}

type TaggedCombiner[T1 any, T2 any] struct {
	*Combiner[T1, T2]
	*TagComponent
}

type AutoInitializer interface {
	Init()
}

type AutoUpdater interface {
	Update()
}

type Toggler interface {
	Toggle()
}
