package pwdhash

import "github.com/allisson/go-pwdhash/argon2"

type config struct {
	current Hasher
}

type Option func(*config)

func defaultConfig() *config {
	return &config{
		current: argon2.Default(),
	}
}

func WithHasher(h Hasher) Option {
	return func(c *config) {
		c.current = h
	}
}
