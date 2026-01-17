package encoding

import (
	"testing"

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

	require.Equal(t, original.Algorithm, parsed.Algorithm)
	require.Equal(t, original.Params, parsed.Params)
	require.Equal(t, original.Salt, parsed.Salt)
	require.Equal(t, original.Hash, parsed.Hash)
}

func TestParse_ErrorScenarios(t *testing.T) {
	const validParams = "m=65536,t=3,p=4"
	const validSalt = "YWJj"
	const validHash = "ZGVm"

	tests := []struct {
		name      string
		input     string
		errSubstr string
	}{
		{
			name:      "missingPrefix",
			input:     "argon2id$v=19$" + validParams + "$" + validSalt + "$" + validHash,
			errSubstr: "invalid PHC string",
		},
		{
			name:      "tooFewParts",
			input:     "$argon2id$",
			errSubstr: "invalid PHC format",
		},
		{
			name:      "missingVersion",
			input:     "$argon2id$version$" + validParams + "$" + validSalt + "$" + validHash,
			errSubstr: "missing version",
		},
		{
			name:      "invalidVersion",
			input:     "$argon2id$v=not-a-number$" + validParams + "$" + validSalt + "$" + validHash,
			errSubstr: "invalid syntax",
		},
		{
			name:      "invalidParam",
			input:     "$argon2id$v=19$m65536$" + validSalt + "$" + validHash,
			errSubstr: "invalid param",
		},
		{
			name:      "badSalt",
			input:     "$argon2id$v=19$" + validParams + "$@@@$" + validHash,
			errSubstr: "illegal base64",
		},
		{
			name:      "badHash",
			input:     "$argon2id$v=19$" + validParams + "$" + validSalt + "$@@@",
			errSubstr: "illegal base64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errSubstr)
		})
	}
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

	require.Contains(t, encoded, "$argon2id$")
	require.Contains(t, encoded, "$v=19$")
	require.Contains(t, encoded, "m=65536")
}
