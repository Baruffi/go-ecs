package ecs

type AnyComponent interface {
}

type Registry struct {
	components map[ComponentId]AnyComponent
}

// NewRegistry - Creates a new registry filling in required initialization parameters
func NewRegistry() *Registry {
	return &Registry{
		components: make(map[ComponentId]AnyComponent),
	}
}

// getComponent - Find or create a component of type C
func getComponent[C ComponentData](r *Registry) (Component[C], bool) {
	for _, v := range r.components {
		if vc, ok := v.(Component[C]); ok {
			return vc, true
		}
	}
	vc := NewComponent[C]()
	return vc, false
}

// getComponentGroup - Find or create a component group of type C
func getComponentGroup[C ComponentData](r *Registry) (ComponentGroup[C], bool) {
	for _, v := range r.components {
		if vg, ok := v.(ComponentGroup[C]); ok {
			return vg, true
		}
	}
	vg := NewComponentGroup[C]()
	return vg, false
}

// Has - Checks if component exists in the registry and is linked to the entity
func Has[C ComponentData](r *Registry, e EntityId) bool {
	if vc, ok := getComponent[C](r); ok {
		_, ok := vc.data[e]
		return ok
	}
	return false
}

// HasGroup - Checks if group exists in the registry and contains the entity
func HasGroup[C ComponentData](r *Registry, e EntityId) bool {
	if vg, ok := getComponentGroup[C](r); ok {
		for _, vc := range vg.members {
			if _, ok := vc.data[e]; ok {
				return true
			}
		}
	}
	return false
}

// Link - Links the component to the respective entity inside the registry
func Link[C ComponentData](r *Registry, e EntityId, c C) TypedComponentId[C] {
	vc, _ := getComponent[C](r)
	vc.data[e] = c
	r.components[vc.id] = vc
	return TypedComponentId[C](vc.id)
}

// Group - Links the component to the respective entity inside the group
func Group[C ComponentData](r *Registry, e EntityId, c ComponentData) TypedComponentId[C] {
	vg, _ := getComponentGroup[C](r)
	vc := NewComponent[C]()
	vc.data[e] = c.(C)
	vg.members[vc.id] = vc
	r.components[vg.id] = vg
	return TypedComponentId[C](vg.id)
}

// Unlink - Unlinks the component type from the respective entity inside the registry (if they were not linked, this is a no-op)
func Unlink[C ComponentData](r *Registry, e EntityId) {
	if vc, ok := getComponent[C](r); ok {
		delete(vc.data, e)
	}
}

// Ungroup - Removes all instances of entity e from group
func Ungroup[C ComponentData](r *Registry, e EntityId) {
	if vg, ok := getComponentGroup[C](r); ok {
		for _, vc := range vg.members {
			delete(vc.data, e)
		}
	}
}

// ClearType - Removes all instances of a component type from the respective the registry
func ClearType[C ComponentData](r *Registry) {
	if vc, ok := getComponent[C](r); ok {
		delete(r.components, vc.id)
	}
}

// ClearGroup - Removes all instances of a component group from the respective the registry
func ClearGroup[C ComponentData](r *Registry) {
	if vg, ok := getComponentGroup[C](r); ok {
		delete(r.components, vg.id)
	}
}

// Clear - Fully clears the respective the registry
func Clear(r *Registry) {
	r.components = make(map[ComponentId]AnyComponent)
}

// View - gets all entities from the component type
func View[C ComponentData](r *Registry) map[EntityId]C {
	if vc, ok := getComponent[C](r); ok {
		return vc.data
	}

	return nil
}

// ViewGroup - gets all entities from the component group
func ViewGroup[C ComponentData](r *Registry) []map[EntityId]C {
	group := make([]map[EntityId]C, 0)
	if vg, ok := getComponentGroup[C](r); ok {
		for _, vc := range vg.members {
			group = append(group, vc.data)
		}
	}

	return group
}

// ViewById - gets all entities from the component type by id
func ViewById[C ComponentData](r *Registry, i TypedComponentId[C]) map[EntityId]C {
	if vc, ok := r.components[ComponentId(i)].(Component[C]); ok {
		return vc.data
	}

	return nil
}

// ViewGroupById - gets all entities from the component group by id
func ViewGroupById[C ComponentData](r *Registry, i TypedComponentId[C]) []map[EntityId]C {
	group := make([]map[EntityId]C, 0)
	if vg, ok := r.components[ComponentId(i)].(ComponentGroup[C]); ok {
		for _, vc := range vg.members {
			group = append(group, vc.data)
		}
	}

	return group
}

// Get - gets specific component data by its type and parent entity id (or returns the default value)
func Get[C ComponentData](r *Registry, e EntityId) (c C, ok bool) {
	vc, ok := getComponent[C](r)
	if ok {
		c, ok = vc.data[e]
	}

	return c, ok
}

// GetGroup - gets specific component data group by its group and parent entity id
func GetGroup[C ComponentData](r *Registry, e EntityId) []C {
	group := make([]C, 0)
	if vg, ok := getComponentGroup[C](r); ok {
		for _, vc := range vg.members {
			group = append(group, vc.data[e])
		}
	}

	return group
}

// GetById - gets specific component data by its component id and parent entity id (or returns the default value)
func GetById[C ComponentData](r *Registry, i TypedComponentId[C], e EntityId) (c C, ok bool) {
	vc, ok := r.components[ComponentId(i)].(Component[C])
	if ok {
		c, ok = vc.data[e]
	}

	return c, ok
}

// GetGroupById - gets specific component data group by its group id and parent entity id
func GetGroupById[C ComponentData](r *Registry, i TypedComponentId[C], e EntityId) []C {
	group := make([]C, 0)
	if vg, ok := r.components[ComponentId(i)].(ComponentGroup[C]); ok {
		for _, vc := range vg.members {
			group = append(group, vc.data[e])
		}
	}

	return group
}
