package encoding

import (
	"encoding/base64"
	"errors"
	"strings"
)

func Parse(s string) (*EncodedHash, error) {
	parts := strings.Split(s, "$")
	if len(parts) < 6 {
		return nil, errors.New("invalid PHC string")
	}

	params := map[string]string{}
	for _, kv := range strings.Split(parts[3], ",") {
		p := strings.SplitN(kv, "=", 2)
		params[p[0]] = p[1]
	}

	salt, _ := base64.RawStdEncoding.DecodeString(parts[4])
	hash, _ := base64.RawStdEncoding.DecodeString(parts[5])

	return &EncodedHash{
		Algorithm: parts[1],
		Params:    params,
		Salt:      salt,
		Hash:      hash,
	}, nil
}
