package ecs

type Stage struct {
	sceneIdx int
	scenes   []*Scene
}

func NewStage(startIdx int, scenes []*Scene) Stage {
	return Stage{
		sceneIdx: startIdx,
		scenes:   scenes,
	}
}

func (s *Stage) NextScene() {
	s.sceneIdx++
	if s.sceneIdx >= len(s.scenes) {
		s.sceneIdx = 0
	}
}

func (s *Stage) GetScene() *Scene {
	return s.scenes[s.sceneIdx]
}
