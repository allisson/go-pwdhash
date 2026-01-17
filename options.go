package pwdhash

import "github.com/allisson/go-pwdhash/argon2"

// config holds PasswordHasher construction settings.
//
// These defaults are internal helpers; prefer With* options for user-facing
// configuration until the API stabilizes.
type config struct {
	current Hasher
}

// Option configures PasswordHasher construction.
type Option func(*config)

// defaultConfig initializes a config with the default Argon2id hasher.
func defaultConfig() *config {
	return &config{
		current: argon2.Default(),
	}
}

// WithHasher overrides the default hashing algorithm.
func WithHasher(h Hasher) Option {
	return func(c *config) {
		c.current = h
	}
}
