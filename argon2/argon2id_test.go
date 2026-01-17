package argon2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestHasher() *Argon2idHasher {
	return &Argon2idHasher{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func TestArgon2idHasher_HashAndVerify(t *testing.T) {
	hasher := newTestHasher()

	encoded, err := hasher.Hash([]byte("password"))
	require.NoError(t, err)

	ok, err := hasher.Verify([]byte("password"), encoded)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = hasher.Verify([]byte("wrong"), encoded)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestArgon2idHasher_NeedsRehashParameterChange(t *testing.T) {
	hasher := newTestHasher()

	encoded, err := hasher.Hash([]byte("password"))
	require.NoError(t, err)

	needs, err := hasher.NeedsRehash(encoded)
	require.NoError(t, err)
	require.False(t, needs)

	mutated := strings.Replace(encoded, "m=65536", "m=32768", 1)

	needs, err = hasher.NeedsRehash(mutated)
	require.NoError(t, err)
	require.True(t, needs)
}

func TestArgon2idHasher_NeedsRehashAlgorithmChange(t *testing.T) {
	hasher := newTestHasher()

	encoded, err := hasher.Hash([]byte("password"))
	require.NoError(t, err)

	mutated := strings.Replace(encoded, "argon2id", "argon2i", 1)

	needs, err := hasher.NeedsRehash(mutated)
	require.NoError(t, err)
	require.True(t, needs)
}
