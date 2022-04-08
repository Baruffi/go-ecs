package ecs

type AnyComponentGroup interface {
	GetId() ComponentGroupId
	Has(EntityId) bool
	Get(EntityId) []ComponentData
	Set(AnyComponent)
	Unset(EntityId)
}

type AnyComponent interface {
	GetId() ComponentId
	Has(EntityId) bool
	Get(EntityId) (ComponentData, bool)
	Set(EntityId, ComponentData)
	Unset(EntityId)
}

type Registry struct {
	componentGroups map[ComponentGroupId]AnyComponentGroup
	components      map[ComponentId]AnyComponent
}

// NewRegistry - Creates a new registry filling in required initialization parameters
func NewRegistry() *Registry {
	return &Registry{
		componentGroups: make(map[ComponentGroupId]AnyComponentGroup),
		components:      make(map[ComponentId]AnyComponent),
	}
}

func (r *Registry) HasComponent(c ComponentId) bool {
	_, ok := r.components[c]
	return ok
}

func (r *Registry) HasGroup(g ComponentGroupId) bool {
	_, ok := r.componentGroups[g]
	return ok
}

func (r *Registry) GetComponent(cid ComponentId) AnyComponent {
	return r.components[cid]
}

func (r *Registry) GetGroup(gid ComponentGroupId) AnyComponentGroup {
	return r.componentGroups[gid]
}

func (r *Registry) SetComponent(c AnyComponent) {
	r.components[c.GetId()] = c
}

func (r *Registry) SetGroup(g AnyComponentGroup) {
	r.componentGroups[g.GetId()] = g
}

func (r *Registry) UnsetComponent(c ComponentId) {
	delete(r.components, c)
}

func (r *Registry) UnsetGroup(g ComponentGroupId) {
	delete(r.componentGroups, g)
}

func (r *Registry) Clear() {
	r.componentGroups = make(map[ComponentGroupId]AnyComponentGroup)
	r.components = make(map[ComponentId]AnyComponent)
}

// getComponent - Find or create a component of type D
func getComponent[D ComponentData](r *Registry) (Component[D], bool) {
	for _, v := range r.components {
		if c, ok := v.(Component[D]); ok {
			return c, true
		}
	}
	c := NewComponent[D]()
	return c, false
}

// getComponentGroup - Find or create a component group of type D
func getComponentGroup[D ComponentData](r *Registry) (ComponentGroup[D], bool) {
	for _, v := range r.componentGroups {
		if g, ok := v.(ComponentGroup[D]); ok {
			return g, true
		}
	}
	g := NewComponentGroup[D]()
	return g, false
}

// Link - Links the component to the respective entity inside the registry
func Link[D ComponentData](r *Registry, e EntityId, d D) ComponentId {
	c, _ := getComponent[D](r)
	c.data[e] = d
	r.SetComponent(c)
	return c.id
}

// Group - Links the component to the respective entity inside the group
func Group[D ComponentData](r *Registry, e EntityId, d ComponentData) (ComponentId, ComponentGroupId) {
	g, _ := getComponentGroup[D](r)
	c := NewComponent[D]()
	c.data[e] = d.(D)
	g.members[c.id] = c
	r.SetGroup(g)
	return c.id, g.id
}

// Has - Checks if component exists in the registry and is linked to the entity
func Has[D ComponentData](r *Registry, e EntityId) bool {
	if c, ok := getComponent[D](r); ok {
		return c.Has(e)
	}
	return false
}

// HasGroup - Checks if group exists in the registry and contains the entity
func HasGroup[D ComponentData](r *Registry, e EntityId) bool {
	if g, ok := getComponentGroup[D](r); ok {
		return g.Has(e)
	}
	return false
}

// Unlink - Unlinks the component type from the respective entity inside the registry (if they were not linked, this is a no-op)
func Unlink[D ComponentData](r *Registry, e EntityId) {
	if c, ok := getComponent[D](r); ok {
		c.Unset(e)
	}
}

// Ungroup - Removes all instances of entity e from group
func Ungroup[D ComponentData](r *Registry, e EntityId) {
	if g, ok := getComponentGroup[D](r); ok {
		g.Unset(e)
	}
}

// ClearType - Removes component type from the respective the registry
func ClearType[D ComponentData](r *Registry) {
	if c, ok := getComponent[D](r); ok {
		r.UnsetComponent(c.id)
	}
}

// ClearGroup - Removes component group from the respective the registry
func ClearGroup[D ComponentData](r *Registry) {
	if g, ok := getComponentGroup[D](r); ok {
		r.UnsetGroup(g.id)
	}
}

// ClearEntity - Removes all instances of an entity from the respective the registry
func ClearEntity(r *Registry, e EntityId) {
	for _, g := range r.componentGroups {
		g.Unset(e)
	}
	for _, c := range r.components {
		c.Unset(e)
	}
}

// View - gets all entities from the component type
func View[D ComponentData](r *Registry) map[EntityId]D {
	if c, ok := getComponent[D](r); ok {
		return c.data
	}

	return nil
}

// ViewGroup - gets all entities from the component group
func ViewGroup[D ComponentData](r *Registry) map[ComponentId]map[EntityId]D {
	group := make(map[ComponentId]map[EntityId]D)
	if g, ok := getComponentGroup[D](r); ok {
		for cid, c := range g.members {
			group[cid] = c.data
		}
	}

	return group
}

// Get - gets specific component data by its type and parent entity id (or returns the default value)
func Get[D ComponentData](r *Registry, e EntityId) (d D, ok bool) {
	c, ok := getComponent[D](r)
	if ok {
		d, ok = c.data[e]
	}

	return d, ok
}

// GetGroup - gets specific component data group by its group and parent entity id
func GetGroup[D ComponentData](r *Registry, e EntityId) map[ComponentId]D {
	group := make(map[ComponentId]D)
	if g, ok := getComponentGroup[D](r); ok {
		for cid, c := range g.members {
			group[cid] = c.data[e]
		}
	}

	return group
}
