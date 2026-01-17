package argon2

import (
	"fmt"
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

func TestParamsForPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		policy int
		want   PolicyParams
		ok     bool
	}{
		{
			name:   "interactive",
			policy: 0,
			want: PolicyParams{
				Memory:      64 * 1024,
				Iterations:  3,
				Parallelism: 4,
			},
			ok: true,
		},
		{
			name:   "moderate",
			policy: 1,
			want: PolicyParams{
				Memory:      128 * 1024,
				Iterations:  4,
				Parallelism: 4,
			},
			ok: true,
		},
		{
			name:   "sensitive",
			policy: 2,
			want: PolicyParams{
				Memory:      256 * 1024,
				Iterations:  5,
				Parallelism: 8,
			},
			ok: true,
		},
		{
			name:   "unknown",
			policy: 3,
			want:   PolicyParams{},
			ok:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			params, err := ParamsForPolicy(tt.policy)
			if tt.ok {
				require.NoError(t, err)
				require.Equal(t, tt.want, params)
				return
			}

			require.Error(t, err)
		})
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

func TestArgon2idHasher_VerifyErrors(t *testing.T) {
	hasher := newTestHasher()

	encoded, err := hasher.Hash([]byte("password"))
	require.NoError(t, err)

	tests := []struct {
		name      string
		mutator   func(string) string
		expectErr bool
		expectOK  bool
	}{
		{
			name: "wrongAlgorithm",
			mutator: func(s string) string {
				return strings.Replace(s, "argon2id", "bcrypt", 1)
			},
			expectErr: false,
			expectOK:  false,
		},
		{
			name: "wrongVersion",
			mutator: func(s string) string {
				return strings.Replace(s, "v=19", "v=27", 1)
			},
			expectErr: true,
		},
		{
			name: "invalidParams",
			mutator: func(s string) string {
				return strings.Replace(s, "m=65536", "m=abc", 1)
			},
			expectErr: true,
		},
		{
			name: "invalidHashLength",
			mutator: func(s string) string {
				return s[:len(s)-4]
			},
			expectErr: false,
			expectOK:  false,
		},
		{
			name: "malformedPHC",
			mutator: func(string) string {
				return "$argon2id$"
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mut := tt.mutator(encoded)
			ok, err := hasher.Verify([]byte("password"), mut)
			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectOK, ok)
		})
	}
}

func TestArgon2idHasher_NeedsRehashParameterChange(t *testing.T) {
	hasher := newTestHasher()

	encoded, err := hasher.Hash([]byte("password"))
	require.NoError(t, err)

	needs, err := hasher.NeedsRehash(encoded)
	require.NoError(t, err)
	require.False(t, needs)

	mutated := strings.Replace(encoded, "m=65536", fmt.Sprintf("m=%d", hasher.Memory/2), 1)

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
