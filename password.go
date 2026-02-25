package pwdhash

import (
	"fmt"

	"github.com/allisson/go-pwdhash/internal/encoding"
)

// PasswordHasher manages password hashing operations via registered algorithms.
type PasswordHasher struct {
	current  Hasher
	registry map[string]Hasher
}

// New constructs a PasswordHasher configured via the provided options.
func New(opts ...Option) (*PasswordHasher, error) {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	reg := make(map[string]Hasher)
	reg[cfg.current.ID()] = cfg.current

	return &PasswordHasher{
		current:  cfg.current,
		registry: reg,
	}, nil
}

// Hash encodes the provided password using the active hasher.
func (p *PasswordHasher) Hash(password []byte) (string, error) {
	return p.current.Hash(password)
}

// Verify checks whether the encoded hash matches the provided password.
func (p *PasswordHasher) Verify(password []byte, encoded string) (bool, error) {
	parsed, err := encoding.Parse(encoded)
	if err != nil {
		return false, err
	}

	hasher, ok := p.registry[parsed.Algorithm]
	if !ok {
		return false, fmt.Errorf("unknown hash algorithm: %s", parsed.Algorithm)
	}

	return hasher.Verify(password, encoded)
}

// NeedsRehash reports whether the encoded hash should be regenerated.
func (p *PasswordHasher) NeedsRehash(encoded string) (bool, error) {
	parsed, err := encoding.Parse(encoded)
	if err != nil {
		return false, err
	}

	hasher, ok := p.registry[parsed.Algorithm]
	if !ok {
		return true, nil
	}

	return hasher.NeedsRehash(encoded)
}
