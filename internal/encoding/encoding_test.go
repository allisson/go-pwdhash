package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_ValidPHCString(t *testing.T) {
	original := &EncodedHash{
		Algorithm: "argon2id",
		Version:   19,
		Params: map[string]string{
			"m": "65536",
			"t": "3",
			"p": "4",
		},
		Salt: []byte("0123456789abcdef"),
		Hash: []byte("abcdefghijklmnopqrstuvwx"),
	}

	encoded := original.String()

	parsed, err := Parse(encoded)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	assert.Equal(t, original.Algorithm, parsed.Algorithm)
	assert.Equal(t, original.Params, parsed.Params)
	assert.Equal(t, original.Salt, parsed.Salt)
	assert.Equal(t, original.Hash, parsed.Hash)
}

func TestParse_InvalidPHCString(t *testing.T) {
	_, err := Parse("$argon2id$")
	require.Error(t, err)
}

func TestEncodedHashStringIncludesMetadata(t *testing.T) {
	enc := EncodedHash{
		Algorithm: "argon2id",
		Version:   19,
		Params: map[string]string{
			"m": "65536",
		},
		Salt: []byte("0123456789abcdef"),
		Hash: []byte("abcdefghijklmnopqrstuvwx"),
	}

	encoded := enc.String()

	assert.Contains(t, encoded, "$argon2id$")
	assert.Contains(t, encoded, "$v=19$")
	assert.Contains(t, encoded, "m=65536")
}
