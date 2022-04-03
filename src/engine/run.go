package engine

type StatefulRunner interface {
	Running() bool
}

type StatefulRunnerFunc func() bool

func (f StatefulRunnerFunc) Running() bool {
	return f()
}

type Runner interface {
	StatefulRunner
	Init()
	Run()
	Step()
}
