# go-pwdhash
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/go-pwdhash)](https://goreportcard.com/report/github.com/allisson/go-pwdhash)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/go-pwdhash)

`go-pwdhash` is a Go-first password hashing helper that embraces the PHC (Password Hashing Competition) format, wraps Argon2id with safe defaults, and surfaces a minimal API for hashing, verification, and upgrades.

## Argon2id Only

pwdhash intentionally supports **Argon2id only**. Algorithms that have already been superseded by Argon2id will not be added, reducing the chance of accidentally selecting outdated primitives. If a superior successor to Argon2id emerges, pwdhash will adopt it behind the same API surface.

## Password Policies

pwdhash ships with opinionated Argon2id policies so applications can select a strength profile without touching raw parameters:

- **Interactive** – user login flows where latency matters most.
- **Moderate** – API keys, service-to-service calls, and privileged automation.
- **Sensitive** – infrastructure secrets, root accounts, and long-lived credentials.

Policies prevent insecure configurations by clamping the underlying Argon2id memory, iteration, and parallelism values to vetted presets.

## Highlights

- **PHC-compliant output** – hashes look like `$argon2id$v=19$...` and parse cleanly across ecosystems.
- **Deterministic upgrade path** – `NeedsRehash` compares stored parameters to the active policy so callers know exactly when to re-encode.
- **Extensible registry** – advanced users may inject tuned Argon2id instances or alternate hashers via the option system.
- **Constant-time verification** – comparisons use helpers under `internal/subtle` to avoid timing leaks.

## Zeroization

pwdhash zeroizes password bytes, salts, and derived keys as soon as they are no longer needed. Due to Go runtime behavior this is a best-effort mitigation, but it significantly shortens the memory exposure window for sensitive material.

## Installation

```bash
go get github.com/allisson/go-pwdhash
```

The module targets Go 1.24, depends on `golang.org/x/crypto`, and uses `stretchr/testify` solely for tests.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/allisson/go-pwdhash"
)

func main() {
    hasher, err := pwdhash.New(
        pwdhash.WithPolicy(pwdhash.PolicyInteractive),
    )
    if err != nil {
        panic(err)
    }

    encoded, err := hasher.Hash([]byte("s3cret"))
    if err != nil {
        panic(err)
    }

    ok, err := hasher.Verify([]byte("s3cret"), encoded)
    if err != nil {
        panic(err)
    }
    fmt.Println("password matches?", ok)

    rehash, err := hasher.NeedsRehash(encoded)
    if err != nil {
        panic(err)
    }
    fmt.Println("should upgrade hash?", rehash)
}
```

## Usage Examples

### Login Handler with Rehashes

```go
func login(ctx context.Context, email, candidate string) error {
    storedHash := loadHashFromDB(ctx, email)

    hasher, err := pwdhash.New(pwdhash.WithPolicy(pwdhash.PolicyModerate))
    if err != nil {
        return fmt.Errorf("pwdhash init: %w", err)
    }

    ok, err := hasher.Verify([]byte(candidate), storedHash)
    if err != nil {
        return fmt.Errorf("pwdhash verify: %w", err)
    }
    if !ok {
        return errInvalidCredentials
    }

    needsUpgrade, err := hasher.NeedsRehash(storedHash)
    if err != nil {
        return fmt.Errorf("pwdhash rehash check: %w", err)
    }
    if needsUpgrade {
        upgraded, err := hasher.Hash([]byte(candidate))
        if err != nil {
            return fmt.Errorf("pwdhash rehash: %w", err)
        }
        persistHash(ctx, email, upgraded)
    }

    return nil
}
```

### Role-Based Policies

```go
type account struct {
    Email string
    Role  string
}

func policyForRole(role string) pwdhash.Policy {
    switch role {
    case "sre", "root":
        return pwdhash.PolicySensitive
    case "service", "api":
        return pwdhash.PolicyModerate
    default:
        return pwdhash.PolicyInteractive
    }
}

func hashForAccount(acct account, password []byte) (string, error) {
    policy := policyForRole(acct.Role)

    hasher, err := pwdhash.New(pwdhash.WithPolicy(policy))
    if err != nil {
        return "", err
    }

    return hasher.Hash(password)
}
```

### Verifying External Hashes

When migrating from another service, you may only have PHC-formatted hashes. You can still apply the pwdhash registry to
validate them while planning a gradual rehash:

```go
func verifyLegacy(hash string, password []byte) (bool, error) {
    hasher, err := pwdhash.New()
    if err != nil {
        return false, err
    }

    ok, err := hasher.Verify(password, hash)
    if err != nil {
        return false, err
    }

    // Decide later whether to rehash with NeedsRehash.
    return ok, nil
}
```

## Configuration

`pwdhash.New` accepts functional options:

- `pwdhash.WithPolicy` selects one of the built-in presets.
- `pwdhash.WithHasher` installs a custom `pwdhash.Hasher` (useful for bespoke Argon2id tuning or for experimenting with future algorithms).

Example of injecting custom parameters:

```go
import (
    "github.com/allisson/go-pwdhash"
    "github.com/allisson/go-pwdhash/argon2"
)

func tunedHasher() (*pwdhash.PasswordHasher, error) {
    argon := &argon2.Argon2idHasher{
        Memory:      128 * 1024,
        Iterations:  4,
        Parallelism: 2,
        SaltLength:  16,
        KeyLength:   32,
    }

    return pwdhash.New(pwdhash.WithHasher(argon))
}
```

## PHC Encoding

pwdhash serializes `encoding.EncodedHash` values using the canonical PHC layout:

```
$argon2id$v=19$m=65536,t=3,p=4$<base64(salt)>$<base64(hash)>
```

The parser validates algorithm identifiers, versions, parameter pairs, and Base64 payloads before handing them to the active hasher.

## Testing & Tooling

```bash
go test ./...
```

Or use the convenience targets:

```bash
make lint   # golangci-lint run -v --fix
make test   # go test -covermode=count -coverprofile=count.out -v ./...
```

## Contributing

1. Fork and clone the repo.
2. Run `go test ./...` (and `make lint`) before sending patches.
3. Keep exports documented, prefer table-driven tests, and stick to `gofmt`/`goimports`.
4. Argon2id is the focus; proposals for new algorithms should include rationale plus end-to-end tests.

## License

MIT licensed. See `LICENSE` for details.
