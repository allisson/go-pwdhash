# Changelog

All notable changes to this project will be documented in this file.

## v0.2.0 - 2026-01-17
- Added policy presets (`PolicyInteractive`, `PolicyModerate`, `PolicySensitive`) plus the `WithPolicy` option for configuring Argon2id without manual parameters.
- Introduced `argon2.ParamsForPolicy` and a shared policy descriptor so the CLI and library reuse the same vetted defaults.
- Hardened `argon2.Argon2idHasher` with parameter validation guards and augmented doc comments across the new public surface.
- Expanded the test suite to cover policy selection and the Argon2 preset helpers, keeping `go test ./...` green.
- Rewrote README to state the Argon2id-only stance and document the policy workflow.

## v0.1.0 - 2026-01-17
- Initial public release of `go-pwdhash` with Argon2id defaults (64MiB memory, 3 iterations, 4 lanes, 16-byte salts, 32-byte keys).
- Introduced PHC-compliant encoder/decoder plus constant-time comparison helpers.
- Added configurable `PasswordHasher` registry with `Hash`, `Verify`, and `NeedsRehash` lifecycle.
- Bundled README, AGENTS guide, and comprehensive test suite (argon2, encoding, subtle, cast) using `go test`/`testify`.
