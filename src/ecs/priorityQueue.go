package ecs

type PriorityLevel int

const (
	Level0 PriorityLevel = iota
	Level1
	Level2
	Level3
	Level4
	Level5
	Level6
	Level7
	Level8
	Level9
)

type PriorityPackage[T any] struct {
	Level   PriorityLevel
	Element []T
}

type PriorityQueue[T any] struct {
	levels [10][]T
}

func (q *PriorityQueue[T]) Insert(level PriorityLevel, elements ...T) {
	pkg := PriorityPackage[T]{level, elements}
	q.Feed(pkg)
}

func (q *PriorityQueue[T]) Feed(pkg PriorityPackage[T]) {
	if q.levels[pkg.Level] == nil {
		q.levels[pkg.Level] = make([]T, 0)
	}
	q.levels[pkg.Level] = append(q.levels[pkg.Level], pkg.Element...)
}

func (q *PriorityQueue[T]) View() []PriorityPackage[T] {
	pkgs := make([]PriorityPackage[T], 0)
	for level, elements := range q.levels {
		pkgs = append(pkgs, PriorityPackage[T]{PriorityLevel(level), elements})
	}
	return pkgs
}

func (q *PriorityQueue[T]) Consume() []PriorityPackage[T] {
	pkgs := q.View()
	q.levels = [10][]T{}
	return pkgs
}

func (q *PriorityQueue[T]) Len() int {
	var totalLen int
	for _, elements := range q.levels {
		totalLen += len(elements)
	}
	return totalLen
}
