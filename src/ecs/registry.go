package ecs

type AnyComponent interface {
}

type Registry struct {
	components map[ComponentId]AnyComponent
}

// NewRegistry - Creates a new registry filling in required initialization parameters
func NewRegistry() *Registry {
	registry := Registry{
		components: make(map[ComponentId]AnyComponent),
	}

	return &registry
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

// Link - Links the component to the respective entity inside the registry
func Link[C ComponentData](r *Registry, e Entity, c C) Component[C] {
	vc, _ := getComponent[C](r)
	vc.data[e.id] = c
	r.components[vc.id] = vc
	return vc
}

// Unlink - Unlinks the component type from the respective entity inside the registry (if they were not linked, this is a no-op)
func Unlink[C ComponentData](r *Registry, e Entity) {
	if vc, ok := getComponent[C](r); ok {
		delete(vc.data, e.id)
	}
}

// ClearType - Removes all instances of a component type from the respective the registry
func ClearType[C ComponentData](r *Registry) {
	if vc, ok := getComponent[C](r); ok {
		delete(r.components, vc.id)
	}
}

// Clear - Clears the respective the registry
func Clear(r *Registry) {
	r.components = make(map[ComponentId]AnyComponent)
}

// Get - gets all entity components from the component type
func Get[C ComponentData](r *Registry) map[EntityId]C {
	if vc, ok := getComponent[C](r); ok {
		return vc.data
	}

	return nil
}

// GetFrom - gets specific component data by its type and parent entity id (or returns the default value)
func GetFrom[C ComponentData](r *Registry, e Entity) (c C, ok bool) {
	vc, ok := getComponent[C](r)
	if ok {
		c = vc.data[e.id]
	}

	return c, ok
}
