package queue

const (
	NOT_DONE HandlerResult = iota
	DONE
)

type HandlerResult int

type Handler[T any] interface {
	Handle(T) HandlerResult
}

type HandlerFunc[T any] func(T) HandlerResult

func (f HandlerFunc[T]) Handle(t T) HandlerResult {
	return f(t)
}
