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

// WithPolicy selects a preset Argon2id configuration for the PasswordHasher.
func WithPolicy(p Policy) Option {
	return func(c *config) {
		params, err := argon2.ParamsForPolicy(int(p))
		if err != nil {
			panic(err) // invalid policy is a programming bug
		}

		c.current = &argon2.Argon2idHasher{
			Memory:      params.Memory,
			Iterations:  params.Iterations,
			Parallelism: params.Parallelism,
			SaltLength:  16,
			KeyLength:   32,
		}
	}
}
