package subtle

import "testing"

func TestConstantTimeCompare(t *testing.T) {
	a := []byte("same-length")
	b := []byte("same-length")

	if !ConstantTimeCompare(a, b) {
		t.Fatalf("expected slices to be equal")
	}

	different := []byte("different")
	if ConstantTimeCompare(a, different) {
		t.Fatalf("expected mismatch to return false")
	}
}
