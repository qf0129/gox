package ptrx

func Ptr[T any](t T) *T {
	return &t
}
