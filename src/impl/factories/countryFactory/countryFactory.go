package countryFactory

import (
	"example.com/v0/src/ecs"
	"example.com/v0/src/impl/systems"
	"github.com/faiface/pixel"
)

func NewFactory(s *ecs.Scene, frame int, position pixel.Vec, orig pixel.Vec, timeLoc string, eventSystem *systems.EventSystem, drawSystem *systems.DrawSystem) ecs.EntityFactory[CountryPrefab] {
	prefab := CountryPrefab{
		frame:      frame,
		position:   position,
		orig:       orig,
		timeLoc:    timeLoc,
		drawSystem: drawSystem,
	}
	return ecs.NewEntityFactory(s, prefab)
}
