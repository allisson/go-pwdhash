package subtle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstantTimeCompare(t *testing.T) {
	a := []byte("same-length")
	b := []byte("same-length")

	require.True(t, ConstantTimeCompare(a, b))

	different := []byte("different")
	require.False(t, ConstantTimeCompare(a, different))
}
