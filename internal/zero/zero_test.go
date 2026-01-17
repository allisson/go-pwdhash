package zero

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesZeroizesSlice(t *testing.T) {
	buf := []byte("supersecret")

	Bytes(buf)

	require.Equal(t, make([]byte, len(buf)), buf)
}

func TestBytesNilSafe(t *testing.T) {
	// Should not panic
	Bytes(nil)
}
