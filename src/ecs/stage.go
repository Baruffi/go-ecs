package ecs

type Stage struct {
	currentScene string
	scenes       map[string]*Scene
}

func NewStage(startScene string, scenes map[string]*Scene) *Stage {
	return &Stage{
		currentScene: startScene,
		scenes:       scenes,
	}
}

func (s *Stage) ChangeScene(scene string) {
	s.currentScene = scene
}

func (s *Stage) GetScene() *Scene {
	return s.scenes[s.currentScene]
}
