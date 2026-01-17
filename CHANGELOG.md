# Changelog

All notable changes to this project will be documented in this file.

## v0.1.0 - 2026-01-17
- Initial public release of `go-pwdhash` with Argon2id defaults (64MiB memory, 3 iterations, 4 lanes, 16-byte salts, 32-byte keys).
- Introduced PHC-compliant encoder/decoder plus constant-time comparison helpers.
- Added configurable `PasswordHasher` registry with `Hash`, `Verify`, and `NeedsRehash` lifecycle.
- Bundled README, AGENTS guide, and comprehensive test suite (argon2, encoding, subtle, cast) using `go test`/`testify`.
