package pwdhash

import "errors"

var (
	// ErrInvalidHash indicates that an encoded hash cannot be parsed.
	ErrInvalidHash = errors.New("invalid encoded hash")
)
