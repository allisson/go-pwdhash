# go-pwdhash

`go-pwdhash` is a Go-first implementation of the Password Hashing Competition (PHC) string format with a batteries-included Argon2id hasher, deterministic upgrade signals, and zero surprises for callers who need predictable password hygiene.

## Highlights

- **Modern Argon2id defaults** – ships with 64MiB memory, 3 iterations, 4 lanes, 16-byte salts, 32-byte keys, and Argon2 v=19 metadata.
- **PHC compliant outputs** – hashes look like `$argon2id$v=19$...` and round-trip cleanly through the built-in parser.
- **Interoperable by design** – encoded hashes verify inside Python's `pwdlib` and equivalent implementations in Rust or C without adapters.
- **Extensible registry** – inject alternative hashers (or tuned Argon2id instances) via options while keeping a single entry point.
- **Deterministic lifecycle** – `Hash`, `Verify`, and `NeedsRehash` expose the minimum API you need to manage password upgrades.
- **Constant-time comparisons** – verification uses `crypto/subtle` helpers to avoid timing leaks.

## Installation

```bash
go get github.com/allisson/go-pwdhash
```

The module targets Go 1.24, depends on `golang.org/x/crypto`, and uses `stretchr/testify` only for tests.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/allisson/go-pwdhash"
)

func main() {
    hasher, err := pwdhash.New()
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

## Configuration

`pwdhash.New` accepts functional options. By default it registers a single Argon2id hasher returned by `argon2.Default()`. To tune parameters, construct the hasher yourself and inject it:

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

To introduce a new algorithm, implement the `pwdhash.Hasher` interface and register it through `WithHasher`. The password hasher keeps an internal registry keyed by `Hasher.ID()`, so mixed fleets of algorithms are possible during migrations.

## PHC Encoding Basics

Internally, the library serializes `encoding.EncodedHash` structures that follow the pattern:

```
$argon2id$v=19$m=65536,t=3,p=4$<base64(salt)>$<base64(hash)>
```

- Parameters are stored verbatim; `NeedsRehash` compares them to the current configuration to decide when to upgrade.
- The parser validates the prefix, version, parameter key/value pairs, and base64 payloads, returning structured errors for callers.

Because the format matches the PHC specification byte-for-byte, it remains compatible with Python, Rust, or C libraries that speak the same dialect.

## Testing & Tooling

Run the suite locally:

```bash
go test ./...
```

Automation-friendly targets live in the Makefile:

```bash
make lint   # golangci-lint run -v --fix
make test   # go test -covermode=count -coverprofile=count.out -v ./...
```

## Contributing

1. Fork and clone the repo.
2. Run `go test ./...` (and `make lint`) before sending patches.
3. Keep exports documented, stick to `gofmt` / `goimports`, and prefer table-driven tests.
4. Discuss larger API changes via issues or draft PRs—Argon2 is the default focus, so new algorithms should include rationale and tests.

## License

MIT licensed. See `LICENSE` for full text.
