package queue

type Numeric interface {
	~int
}

// Max - Return the larger number between a and b
func Max[T Numeric](a T, b T) T {
	if a > b {
		return a
	}
	return b
}
