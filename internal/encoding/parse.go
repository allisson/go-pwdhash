// Package encoding handles serialization and parsing of PHC strings.
package encoding

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

// Parse decodes a PHC-formatted string into an EncodedHash structure.
func Parse(s string) (*EncodedHash, error) {
	if !strings.HasPrefix(s, "$") {
		return nil, errors.New("invalid PHC string")
	}

	parts := strings.Split(s, "$")
	if len(parts) < 6 {
		return nil, errors.New("invalid PHC format")
	}

	algo := parts[1]

	// v=19
	versionPart := parts[2]
	if !strings.HasPrefix(versionPart, "v=") {
		return nil, errors.New("missing version")
	}
	version, err := strconv.Atoi(strings.TrimPrefix(versionPart, "v="))
	if err != nil {
		return nil, err
	}

	params := map[string]string{}
	for _, kv := range strings.Split(parts[3], ",") {
		p := strings.SplitN(kv, "=", 2)
		if len(p) != 2 {
			return nil, errors.New("invalid param")
		}
		params[p[0]] = p[1]
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, err
	}

	return &EncodedHash{
		Algorithm: algo,
		Version:   version,
		Params:    params,
		Salt:      salt,
		Hash:      hash,
	}, nil
}
