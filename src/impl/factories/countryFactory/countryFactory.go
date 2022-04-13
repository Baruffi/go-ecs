package countryFactory

import (
	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/managers"
	"github.com/faiface/pixel"
)

func NewFactory(s *ecs.Scene, frame int, position pixel.Vec, orig pixel.Vec, timeLoc string, eventManager *managers.EventManager, drawerManager *managers.DrawerManager) ecs.EntityFactory[CountryPrefab] {
	prefab := CountryPrefab{
		frame:         frame,
		position:      position,
		orig:          orig,
		timeLoc:       timeLoc,
		drawerManager: drawerManager,
	}
	return ecs.NewEntityFactory(s, prefab)
}
