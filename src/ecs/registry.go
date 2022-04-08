package ecs

type AnyComponentGroup interface {
	GetId() ComponentGroupId
	Has(EntityId) bool
	Get(EntityId) []ComponentData
	Set(AnyComponent)
	Unset(ComponentId)
	UnsetEntity(EntityId)
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
func Link[D ComponentData](r *Registry, e EntityId, d D) TypedComponentId[D] {
	c, _ := getComponent[D](r)
	c.data[e] = d
	r.SetComponent(c)
	return TypedComponentId[D](c.id)
}

// Group - Links the component to the respective entity inside the group
func Group[D ComponentData](r *Registry, e EntityId, d ComponentData) (TypedComponentId[D], TypedComponentGroupId[D]) {
	g, _ := getComponentGroup[D](r)
	c := NewComponent[D]()
	c.data[e] = d.(D)
	g.members[c.id] = c
	r.SetGroup(g)
	return TypedComponentId[D](c.id), TypedComponentGroupId[D](g.id)
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

// HasById - Checks if component exists in the registry by id and is linked to the entity
func HasById[D ComponentData](r *Registry, i TypedComponentId[D], e EntityId) bool {
	if c, ok := r.components[ComponentId(i)]; ok {
		if c.Has(e) {
			return true
		}
	}
	return false
}

// HasGroupById - Checks if group exists in the registry by id and contains the entity
func HasGroupById[D ComponentData](r *Registry, i TypedComponentGroupId[D], e EntityId) bool {
	if g, ok := r.componentGroups[ComponentGroupId(i)]; ok {
		if g.Has(e) {
			return true
		}
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
		g.UnsetEntity(e)
	}
}

// UnlinkById - Unlinks the component of id i from the respective entity inside the registry (if they were not linked, this is a no-op)
func UnlinkById[D ComponentData](r *Registry, i TypedComponentId[D], e EntityId) {
	r.components[ComponentId(i)].Unset(e)
}

// Ungroup - Removes all instances of entity e from group of id i
func UngroupById[D ComponentData](r *Registry, i TypedComponentGroupId[D], e EntityId) {
	r.componentGroups[ComponentGroupId(i)].UnsetEntity(e)
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

// ClearTypeById - Removes component of id i from the respective registry
func ClearTypeById[D ComponentData](r *Registry, i TypedComponentId[D]) {
	r.UnsetComponent(ComponentId(i))
}

// ClearGroupById - Removes group of id i from the respective registry
func ClearGroupById[D ComponentData](r *Registry, i TypedComponentGroupId[D]) {
	r.UnsetGroup(ComponentGroupId(i))
}

// ClearEntity - Removes all instances of an entity from the respective the registry
func ClearEntity(r *Registry, e EntityId) {
	for _, c := range r.componentGroups {
		c.UnsetEntity(e)
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
func ViewGroup[D ComponentData](r *Registry) []map[EntityId]D {
	group := make([]map[EntityId]D, 0)
	if g, ok := getComponentGroup[D](r); ok {
		for _, c := range g.members {
			group = append(group, c.data)
		}
	}

	return group
}

// ViewById - gets all entities from the component type by id
func ViewById[D ComponentData](r *Registry, i TypedComponentId[D]) map[EntityId]D {
	if c, ok := r.components[ComponentId(i)].(Component[D]); ok {
		return c.data
	}

	return nil
}

// ViewGroupById - gets all entities from the component group by id
func ViewGroupById[D ComponentData](r *Registry, i TypedComponentGroupId[D]) []map[EntityId]D {
	group := make([]map[EntityId]D, 0)
	if g, ok := r.componentGroups[ComponentGroupId(i)].(ComponentGroup[D]); ok {
		for _, c := range g.members {
			group = append(group, c.data)
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
func GetGroup[D ComponentData](r *Registry, e EntityId) []D {
	group := make([]D, 0)
	if g, ok := getComponentGroup[D](r); ok {
		for _, c := range g.members {
			group = append(group, c.data[e])
		}
	}

	return group
}

// GetById - gets specific component data by its component id and parent entity id (or returns the default value)
func GetById[D ComponentData](r *Registry, i TypedComponentId[D], e EntityId) (d D, ok bool) {
	c, ok := r.components[ComponentId(i)].(Component[D])
	if ok {
		d, ok = c.data[e]
	}

	return d, ok
}

// GetGroupById - gets specific component data group by its group id and parent entity id
func GetGroupById[D ComponentData](r *Registry, i TypedComponentGroupId[D], e EntityId) []D {
	group := make([]D, 0)
	if g, ok := r.componentGroups[ComponentGroupId(i)].(ComponentGroup[D]); ok {
		for _, c := range g.members {
			group = append(group, c.data[e])
		}
	}

	return group
}

// Ids are not worth it right now
// GetFromGroup - gets specific component data by its component id and parent entity id from inside a group (or returns the default value)
func GetFromGroup[D ComponentData](r *Registry, gi TypedComponentGroupId[D], ci TypedComponentId[D], e EntityId) (d D, ok bool) {
	var c Component[D]
	g, ok := r.componentGroups[ComponentGroupId(gi)].(ComponentGroup[D])
	if ok {
		c, ok = g.members[ComponentId(ci)]
		if ok {
			d, ok = c.data[e]
		}
	}

	return d, ok
}
