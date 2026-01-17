package encoding

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type EncodedHash struct {
	Algorithm string
	Version   int
	Params    map[string]string
	Salt      []byte
	Hash      []byte
}

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
