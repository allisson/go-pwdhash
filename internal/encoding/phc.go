package encoding

import (
	"encoding/base64"
	"fmt"
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

// String renders the hash in PHC string format.
func (e EncodedHash) String() string {
	params := []string{}
	for k, v := range e.Params {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
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
