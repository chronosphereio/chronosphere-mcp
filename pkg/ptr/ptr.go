// Package ptr provides a utility to create and dereference pointers.
package ptr

// To creates a pointer to the value of type T.
func To[T any](v T) *T {
	return &v
}

// Value dereferences a pointer of type T. If the pointer is nil, it returns the zero value of T.
func Value[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
