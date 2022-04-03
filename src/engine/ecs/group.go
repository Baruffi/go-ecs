package ecs

type ComponentGroup[C ComponentData] struct {
	id      ComponentId
	members map[ComponentId]Component[C]
}

// NewComponentGroup - Creates a new ComponentGroup filling in required initialization parameters
func NewComponentGroup[C ComponentData]() ComponentGroup[C] {
	return ComponentGroup[C]{
		id:      ComponentId(GenerateId()),
		members: make(map[ComponentId]Component[C]),
	}
}
