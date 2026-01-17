package zero

// Bytes overwrites the provided buffer with zeros before it goes out of scope.
//
// Best-effort: Go does not guarantee that the compiler will keep these stores,
// but this pattern is widely used to reduce the lifetime of secrets in memory.
func Bytes(b []byte) {
	if b == nil {
		return
	}
	for i := range b {
		b[i] = 0
	}
}
