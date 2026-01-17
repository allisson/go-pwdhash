// Package encoding handles encoding and parsing of PHC strings.
package encoding

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
)

// EncodedHash captures a PHC-formatted hash and associated metadata.
type EncodedHash struct {
	Algorithm string
	Version   int
	Params    map[string]string
	Salt      []byte
	Hash      []byte
}

// String renders the hash in PHC string format with deterministic parameter ordering.
func (e EncodedHash) String() string {
	keys := make([]string, 0, len(e.Params))
	for k := range e.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	params := make([]string, 0, len(keys))
	for _, k := range keys {
		params = append(params, fmt.Sprintf("%s=%s", k, e.Params[k]))
	}

	return fmt.Sprintf(
		"$%s$v=%d$%s$%s$%s",
		e.Algorithm,
		e.Version,
		strings.Join(params, ","),
		base64.RawStdEncoding.EncodeToString(e.Salt),
		base64.RawStdEncoding.EncodeToString(e.Hash),
	)
}
