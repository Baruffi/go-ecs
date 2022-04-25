package ecs

type Stage struct {
	Clock
	currentIdentifier string
	identifiers       map[string]int
	scenes            []*Scene
	stagehands        []*Stagehand
	active            bool
}

func NewStage() *Stage {
	stage := &Stage{
		identifiers: make(map[string]int),
		scenes:      make([]*Scene, 0),
		stagehands:  make([]*Stagehand, 0),
	}

	return stage
}

func (stage *Stage) Include(identifier string, scene *Scene, updater Updater[*Stage]) *Stage {
	stage.identifiers[identifier] = len(stage.scenes)
	stage.scenes = append(stage.scenes, scene)
	stage.stagehands = append(stage.stagehands, &Stagehand{
		stage:   stage,
		updater: updater,
	})

	return stage
}

func (stage *Stage) Start(identifier string) {
	stage.currentIdentifier = identifier
	stage.active = true
	stage.Init()
}

func (stage *Stage) End() {
	stage.active = false
}

func (stage *Stage) IsActive() bool {
	return stage.active
}

func (stage *Stage) GetScene() *Scene {
	return stage.scenes[stage.identifiers[stage.currentIdentifier]]
}

func (stage *Stage) GetStagehand() *Stagehand {
	return stage.stagehands[stage.identifiers[stage.currentIdentifier]]
}
