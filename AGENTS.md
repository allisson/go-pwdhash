# AGENTS GUIDE

## Mission
- Deliver high-signal, minimally invasive contributions to `github.com/allisson/go-pwdhash`.
- Preserve the Argon2id defaults and make behavior explicit in both docs and code.
- Favor clarity, determinism, and verifiability over cleverness.
- Treat this document as the single source of truth for workflow expectations.

## Repository Snapshot
- Go module: `github.com/allisson/go-pwdhash`.
- Go toolchain: 1.24.x (see `go.mod`). Use `go env GOROOT` if the path matters.
- Entrypoints: `password.go`, `options.go`, `hasher.go`, `errors.go`.
- Argon2 implementation: `argon2/argon2id.go` (+ tests in `argon2/argon2id_test.go`).
- PHC helpers: `internal/encoding` (`parse.go`, `phc.go`).
- Constant-time helpers: `internal/subtle`.
- Casting helpers: `internal/cast`.
- Tests rely on `go test`; Makefile wraps lint/test workflows.
- No Cursor (`.cursor/rules`, `.cursorrules`) or Copilot (`.github/copilot-instructions.md`) rule files exist as of this revision; keep checking when pulling.

## Toolchain & Environment
- Use the system `go` binary (Go 1.24) with module-aware mode.
- `golangci-lint` must be installed (`brew install golangci-lint` if missing).
- Always run `gofmt` and `goimports` (or `golines`) before committing.
- Randomness must come from `crypto/rand`; never fall back to `math/rand`.
- External deps are limited to `golang.org/x/crypto` (+ transitives like `x/sys`) and `stretchr/testify` for tests.
- Respect PHC format semantics (per RFC 9106) when touching encoding code.

## Build, Lint & Test Commands
- `go build ./...` – compile all packages.
- `go vet ./...` – lightweight static checks when editing low-level code.
- `make lint` – runs `golangci-lint run -v --fix` (preferred entry point).
- `make test` – runs `go test -covermode=count -coverprofile=count.out -v ./...`.
- `go test ./... -race` – use when dealing with shared state or concurrency.
- `go test ./argon2 -run TestArgon2idHasher` – target a single test (substitute regex as needed).
- `go test ./internal/encoding -run '^TestParse_ValidPHCString$' -count=1` – template for exact single-test reruns.
- `go test -bench=. ./argon2` – benchmark template (add when needed).
- `go tool cover -html=count.out` – visualize coverage after `make test`.

## Testing Playbook
- Prefer table-driven tests with explicit inputs/expectations.
- Default to byte slice literals for salts/hashes to avoid decode helpers.
- Avoid seeded randomness; rely on deterministic inputs or injectable hooks.
- Use `t.Helper()` for shared helper functions inside `_test.go` files.
- Naming: `TestPackage_Feature_Condition`.
- Single-test reruns: `go test ./pkg -run '^TestName$' -count=1`.
- Race reproduction: `GOEXPERIMENT=arenas go test -run TestName ./pkg` only if arenas are already used.
- Clean coverage artifacts (`count.out`) if tooling changes them.
- Place benchmarks next to tests in the same `_test.go` file.
- Avoid build tags unless absolutely necessary; prefer runtime switches/options.

## Linting & Formatting
- Run `golangci-lint` with `--fix` (via `make lint`); commit auto-fixes.
- Always `gofmt -w` and `goimports -w` touched files.
- Import order: stdlib, blank line, third-party, blank line, local (`github.com/allisson/go-pwdhash/...`).
- Keep lines ≤ 100 chars when practical; wrap long argument lists vertically.
- Preserve trailing commas in multi-line literals for minimal diffs.
- No commented-out code or stray debug prints.
- Doc comments for exported symbols must start with the symbol name and be full sentences.
- Use `go env GOPATH` if you need to inspect module cache paths.

## Code Style: Language Constructs
- Keep interfaces minimal (e.g., `Hasher` is intentionally small).
- Prefer value semantics unless mutation after construction is required.
- Order struct fields by runtime importance (critical knobs first).
- Avoid global mutable state; inject dependencies through options.
- Use `(value, error)` returns; panic only in fatal init paths.
- Wrap errors with context: `fmt.Errorf("argon2id encode: %w", err)`.
- Keep exported errors in `errors.go` (`ErrSomething`) and use `errors.Is`.
- Use `const` for algorithm IDs/default parameters once stable.
- Copy byte slices before returning if callers might mutate them.
- Reuse `internal/subtle` for constant-time operations.

## Code Style: Imports & Organization
- Maintain tight package surfaces; unexport helpers unless necessary.
- Avoid circular dependencies; place shared helpers under `internal/...`.
- Do not re-export internal symbols—wire functionality via options instead.
- Mirror the `argon2` package structure when adding new hashers (constructor + Hash/Verify/NeedsRehash trio).
- Keep PHC parsing strict about syntax while tolerant of parameter order.
- Follow package naming conventions (`pwdhash`, `argon2`, `encoding`, etc.).
- Top-level package comments should summarize purpose in 2–3 sentences.
- Group option helpers in `options.go`; avoid scattered `WithX` functions.
- Ensure deterministic serialization (sort map keys before joining strings).
- Keep constants near their usage sites unless shared widely.

## Error Handling & Logging
- Propagate errors upward; let callers decide how to handle them.
- Never log or print secrets (password bytes, salts, entire hashes).
- Use sentinel errors like `ErrInvalidHash` for invalid user input.
- Be explicit in error strings: mention both action and algorithm.
- In `Verify`, differentiate parse errors from mismatched algorithms.
- `NeedsRehash` must return `true` for unknown algorithms.
- Validate external params (memory, iterations, parallelism) before use; clamp if necessary.
- Keep parser tolerant of param ordering but strict about missing keys or malformed values.
- Bubble up salt generation failures immediately.
- Reserve `errors.Join` for multiple independent errors; otherwise wrap singly.

## Dependency & Module Guidance
- Add deps via `go get module@version`, then run `go mod tidy` and inspect diffs.
- Keep the dependency graph lean; avoid optional crypto libs unless justified.
- Vendor only if build systems require it (default to modules).
- Update `golang.org/x/crypto` deliberately; rerun tests and `go vet` afterward.
- Document new dependency rationales (PR descriptions + AGENTS.md if workflow changes).
- Respect semantic import versioning; avoid pseudo-versions unless required for security.
- Use `replace` directives sparingly and remove before merging unless essential.
- Keep generated files out of version control unless reproducible/required.
- Ensure `go mod tidy` leaves the tree clean.
- Review `go.sum` additions for suspicious modules.

## Working With Password Hashers
- `PasswordHasher` maintains a registry keyed by `Hasher.ID()`.
- New hashers must implement `Hasher` (`ID`, `Hash`, `Verify`, `NeedsRehash`).
- `Hash` must return PHC-formatted strings; centralize encoding via `internal/encoding`.
- `Verify` must parse once, reuse parsed data, and perform constant-time comparisons.
- `NeedsRehash` checks algorithm + parameter drift; treat version changes as rehash triggers.
- Provide option helpers (`WithHasher`) for injecting alternative defaults.
- Never expose raw salts; keep them internal.
- Update README examples when user-visible behavior changes.
- Maintain constant-time comparisons via `internal/subtle`.
- Use `encoding.EncodedHash` when adding metadata to keep PHC compatibility.

## Agent Workflow Expectations
- Start sessions with `git status -sb` and `git diff` to understand the working tree.
- Keep branches focused; avoid mixing refactors with behavioral changes.
- Update this AGENTS guide whenever tooling or style rules evolve.
- Reference paths using backticked relative paths in responses (e.g., `password.go`).
- Default validation: `make lint && make test` when scope is unclear.
- Use the todo tool for multi-step tasks; keep one item `in_progress` at a time.
- Summaries must cover rationale, touched files, and validation steps.
- Tests are mandatory for bug fixes/features; call out any gaps explicitly.
- Use descriptive branch names (`feature/new-hasher`, `fix/phc-parser`, etc.).
- Document single-test commands in PRs when reproducing bugs.

## Cursor/Copilot Rules
- No Cursor rules or Copilot instruction files exist in this repo as of the latest commit.
- Re-scan `.cursor/` and `.github/` on new branches; update this section if rules appear.
- Until such files exist, follow the conventions outlined above.

## Quick Reference: Single-Test Commands
- `go test ./argon2 -run '^TestArgon2idHasher_HashAndVerify$' -count=1`
- `go test ./internal/encoding -run '^TestParse_ErrorScenarios$' -count=1`
- `go test ./... -run '^TestPackage_Feature_Condition$' -count=1` (generic template)

Keep this guide ~150 lines; update responsibly when workflows or tooling change.
