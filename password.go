package pwdhash

import (
	"fmt"

	"github.com/allisson/go-pwdhash/internal/encoding"
)

type PasswordHasher struct {
	current  Hasher
	registry map[string]Hasher
}

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

func (p *PasswordHasher) Hash(password []byte) (string, error) {
	return p.current.Hash(password)
}

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
