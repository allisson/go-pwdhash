// Package pwdhash is a Go-first password hashing helper that embraces the PHC
// (Password Hashing Competition) format.
//
// It wraps Argon2id with safe defaults and surfaces a minimal API for hashing,
// verification, and upgrades. pwdhash intentionally supports Argon2id only,
// reducing the chance of accidentally selecting outdated primitives. If a
// superior successor emerges, pwdhash will adopt it behind the same API
// surface.
//
// pwdhash ships with opinionated Argon2id policies so applications can select
// a strength profile without touching raw parameters (Interactive, Moderate,
// and Sensitive).
package pwdhash
