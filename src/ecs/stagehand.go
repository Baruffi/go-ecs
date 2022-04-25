package ecs

type Updater[T any] interface {
	Update(T, float64)
}

type UpdaterFunc[T any] func(T, float64)

func (f UpdaterFunc[T]) Update(target T, dt float64) {
	f(target, dt)
}

type Stagehand struct {
	stage   *Stage
	updater Updater[*Stage]
}

func (hand *Stagehand) Update() {
	hand.updater.Update(hand.stage, hand.stage.Dt())
}
