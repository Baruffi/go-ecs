package components

import "github.com/faiface/pixel"

type Layer int

const (
	Layer0 Layer = iota
	Layer1
	Layer2
	Layer3
	Layer4
	Layer5
	Layer6
	Layer7
	Layer8
	Layer9
)

type Drawer interface {
	GetLayer() Layer
	Draw(pixel.Target)
}

type SortableDrawerQueue struct {
	Drawers []Drawer
}

func (q *SortableDrawerQueue) Init() {
	q.Drawers = make([]Drawer, 0)
}

func (q *SortableDrawerQueue) Add(d Drawer) {
	q.Drawers = append(q.Drawers, d)
}

func (q *SortableDrawerQueue) Len() int {
	return len(q.Drawers)
}

func (q *SortableDrawerQueue) Less(i, j int) bool {
	return q.Drawers[i].GetLayer() > q.Drawers[j].GetLayer()
}

func (q *SortableDrawerQueue) Swap(i, j int) {
	temp := q.Drawers[i]
	q.Drawers[i] = q.Drawers[j]
	q.Drawers[j] = temp
}
