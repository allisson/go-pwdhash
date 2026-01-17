# AGENTS GUIDE

## Mission
- Provide high-signal, minimally-invasive contributions to `github.com/allisson/go-pwdhash`.
- Preserve current Argon2id defaults while making behavior explicit in docs and code.
- Prefer clarity over cleverness; prioritize deterministic, testable changes.
- Treat this file as the single source of truth for agent workflow expectations.

## Repository Snapshot
- Go module name: `github.com/allisson/go-pwdhash`.
- Target Go toolchain: 1.24.0 (from go.mod); run `go env GOROOT` if you need its exact path.
- Main entry points live in `password.go`, `options.go`, `hasher.go`, and `errors.go`.
- Argon2-specific logic sits in `argon2/argon2id.go`.
- PHC string encode/decode helpers live under `internal/encoding`.
- Constant-time helpers live under `internal/subtle`.
- Tests currently rely on `go test`; there are no third-party test runners.
- `Makefile` wraps linting (`golangci-lint`) and coverage-enabled tests.
- There are no Cursor (.cursor/rules, .cursorrules) or Copilot (.github/copilot-instructions.md) rule files present; keep verifying when the repo updates.

## Toolchain & Environment
- Use the system `go` binary that matches Go 1.24.0 features (e.g., generics, toolchain aware builds).
- `golangci-lint` must be available on PATH; install via `brew install golangci-lint` if missing.
- Prefer `gofmt` and `goimports` (or `golines` for wrapping) before committing changes.
- Randomness relies on `crypto/rand`; do not swap to math/rand even in tests.
- External dependencies are limited to `golang.org/x/crypto` and transitively `golang.org/x/sys`.
- Respect the PHC string format when touching encoding code; follow RFC 9106 guidance if unsure.

## Build & Verification Commands
- `go build ./...` compiles the entire module; use before pushing non-test changes.
- `go build ./cmd/...` would compile CLI entry points if they are ever added; keep the pattern consistent.
- `go vet ./...` catches common mistakes; run when touching low-level code.
- `make lint` executes `golangci-lint run -v --fix`; prefer the Makefile target for reproducibility.
- `make test` executes `go test -covermode=count -coverprofile=count.out -v ./...`.
- `go test ./... -race` is recommended when dealing with concurrency or shared buffers.
- `go test ./argon2 -run TestArgon2idHasher` runs a single test or subtest by regex; adjust package path as needed.
- `go test ./internal/encoding -run TestParse` is the pattern for targeted PHC parser tests.
- `go test -bench=. ./argon2` is the template once benchmarks are added; keep the structure consistent.
- `go tool cover -html=count.out` visualizes coverage after `make test`.

## Testing Playbook
- Default to table-driven tests; keep inputs and expected outputs explicit.
- Favor byte slice literals over hex strings for salts and hashes to reduce decode steps.
- Avoid seeding randomness in tests; instead, inject deterministic salts via helper hooks if new code requires it.
- Use `t.Helper()` in shared test utilities once added.
- Keep test names descriptive: `TestPackage_Feature_Condition`.
- Targeted single-test workflow: `go test ./package -run '^TestName$' -count=1`.
- When reproducing race bugs, run `GOEXPERIMENT=arenas go test -run TestName ./pkg` only if the code uses arenas; otherwise keep env clean.
- Clean artifacts such as `count.out` after CI if build scripts change them.
- Add benchmarks in `_test.go` files alongside unit tests to keep behavior colocated.
- Use build tags sparingly; prefer runtime switches or options for hasher differences.

## Linting & Formatting
- `golangci-lint` should run with `--fix`; commit the auto-fixes.
- Enable linters such as `errcheck`, `gosimple`, `staticcheck`, and `revive` via `.golangci.yml` when/if it is introduced.
- Run `gofmt -w` and `goimports -w` on touched files; do not rely solely on IDE format-on-save because CI expects deterministic formatting.
- Keep import groups ordered: stdlib, third-party, then local (`github.com/allisson/go-pwdhash/...`).
- Avoid blank lines within the same import group; use a single blank line between groups.
- Limit lines to ≤ 100 chars where practical; wrap function arguments vertically when readability benefits.
- Preserve trailing commas in multi-line literals to keep diffs small.
- Do not commit commented-out code or debugging prints.
- Document exported functions with complete sentences; start doc comments with the symbol name.
- Use `go env GOPATH` to confirm workspace-specific tool installs when debugging lint setup issues.

## Code Style: Language Constructs
- Interfaces should stay minimal like `Hasher` (only methods that every implementation truly needs).
- Prefer value semantics for config structs unless mutation is required post-initialization.
- Keep struct fields ordered by importance: runtime-critical knobs first, book-keeping last.
- Avoid global mutable state; additional hashers should be injected via options or registries.
- Return `(value, error)` and let callers branch; avoid panic except in truly fatal init paths.
- Wrap errors via `fmt.Errorf("context: %w", err)` when adding new failure paths.
- Keep exported errors in `errors.go`; name them `ErrSomething` and use `errors.Is` for comparisons.
- Use `const` for algorithm IDs and default parameters once they stabilize.
- Byte slices should be copied before returning if the caller might mutate them.
- Emphasize constant-time comparisons when dealing with secrets; reuse `internal/subtle` utilities instead of re-implementing logic.

## Code Style: Imports & Organization
- Each package should expose a tight public surface; unexport helper functions unless needed externally.
- Avoid circular dependencies; keep helpers in `internal/...` when they are not part of the API contract.
- Re-exporting symbols from internal packages is discouraged; prefer wiring functions through options.
- For new hashers, mirror the `argon2` package layout (constructor, parameters, Hash/Verify/NeedsRehash trio).
- Keep PHC parsing and formatting strict; validate lengths and parameter presence when expanding functionality.
- When adding new files, follow the existing package naming conventions (`pwdhash`, `argon2`, `encoding`).
- Top-level doc comments should explain the package purpose in 2–3 sentences.
- Group related option helpers in `options.go`; avoid scattering `WithSomething` functions.
- Maintain deterministic map iteration when serializing (e.g., sort keys before joining) if output stability matters.
- Keep constants near their usage sites unless they become shared defaults.

## Error Handling & Logging
- Propagate errors upward; let callers decide whether to log or mask them.
- Avoid logging secrets; never print password bytes or full hashes.
- Use sentinel errors for invalid input (e.g., `ErrInvalidHash`).
- When returning fmt errors, mention both the action and the failing algorithm ("argon2id encode: %w").
- In Verify flows, differentiate between parse errors and algorithm mismatches to aid debugging.
- `NeedsRehash` should treat unknown algorithms as `true` (already implemented); maintain that contract for new hashers.
- Validate external parameters (memory, iterations, parallelism) before using them; clamp to safe ranges.
- Keep PHC parser tolerant of parameter ordering but strict about missing values.
- Ensure random salt generation errors bubble up immediately.
- Prefer `errors.Join` only when multiple independent failures need reporting; otherwise keep single errors.

## Dependency & Module Guidance
- Use `go get module@version` to add dependencies; run `go mod tidy` afterward and review diff.
- Keep the dependency tree small; avoid optional crypto libraries unless justified.
- Vendor dependencies only if build systems require it; default to modules.
- Update `golang.org/x/crypto` deliberately and re-run tests plus `go vet`.
- Document any new dependency rationale inside PR descriptions and in AGENTS.md if it affects workflow.
- Respect semantic import versioning; do not use pseudo-versions unless necessary for security patches.
- Use replace directives sparingly and remove them before merging unless absolutely required.
- Keep generated files (if any) out of version control unless reproducible and essential.
- Validate that `go mod tidy` leaves a clean state before committing.
- Cross-check `go.sum` additions to ensure no malicious transitive modules slip in.

## Working With Password Hashers
- `PasswordHasher` holds a registry keyed by algorithm IDs; when adding new hashers, register them during construction.
- New hashers should implement `Hasher` and mimic the Argon2id method signatures.
- `Hash` should return PHC-formatted strings; keep encoding centralized in `internal/encoding`.
- `Verify` must parse the PHC string once and reuse parsed data for comparisons.
- `NeedsRehash` should check both algorithm ID and parameter drift; treat version upgrades as rehash triggers.
- Provide option helpers (e.g., `WithHasher`) for injecting alternative default hashers.
- Do not expose `[ ]byte` salts outside the package; keep them internal for security.
- Add doc examples in README when user-visible behavior changes.
- Maintain constant-time comparisons via `internal/subtle.ConstantTimeCompare`.
- Consider hooking into `encoding.EncodedHash` when adding metadata; keep compatibility with other PHC consumers.

## Agent Workflow Expectations
- Start by running `git status` and `git diff` to stay aware of pre-existing local changes.
- Keep branches focused; avoid mixing refactors with behavioral changes.
- Update this AGENTS.md whenever tooling, style, or rules evolve.
- Reference file paths with backticked relative paths in PR or chat responses (e.g., `password.go`).
- When CI requirements are unknown, default to `make lint && make test` locally.
- For large changes, create a todo list (via provided tooling) to track sub-tasks.
- Summaries to users should include rationale, touched files, and validation steps.
- Tests are mandatory for bug fixes and new features; explain any gaps explicitly.
- Prefer descriptive branch names like `feature/new-hasher` or `fix/phc-parser`.
- Document single-test commands in PR descriptions when reproducing specific failures.

## Cursor/Copilot Rules
- No Cursor rules (.cursor/rules, .cursorrules) are present as of this writing.
- No Copilot instructions (.github/copilot-instructions.md) are present as of this writing.
- Check for new rule files when pulling updates and update this section if they appear.
- If rules are added later, summarize their impact here so future agents stay aligned.
- Until such files exist, default to the conventions documented above.
