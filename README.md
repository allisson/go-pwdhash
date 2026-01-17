# pwdhash

`pwdhash` is a Go library for producing PHC-compatible password hashes with sensible Argon2id defaults, explicit configuration, and deterministic rehash guidance.

## Why pwdhash?

- **Modern defaults** – ships with Argon2id tuned to `64MiB` memory, `3` iterations, `4` lanes, `16` byte salts, and `32` byte keys.
- **PHC compliance** – encodes outputs as `$argon2id$v=19$...` strings so you can interoperate with other runtimes.
- **Extensible registry** – swap or add hashers through options without touching call sites.
- **Clear lifecycle** – `Hash`, `Verify`, and `NeedsRehash` expose the minimal API you need for password management.
- **Constant-time verification** – comparisons run through `crypto/subtle` helpers to avoid timing leaks.

## Installation

```bash
go get github.com/allisson/go-pwdhash
```

The module targets Go 1.24 and depends only on `golang.org/x/crypto` plus `testify` for tests.

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

    needsRehash, err := hasher.NeedsRehash(encoded)
    if err != nil {
        panic(err)
    }
    fmt.Println("should upgrade hash?", needsRehash)
}
```

## Configuration

`pwdhash.New` accepts functional options. By default it registers a single Argon2id hasher returned by `argon2.Default()`. To customize parameters, call the constructor directly and inject it:

```go
import "github.com/allisson/go-pwdhash/argon2"

func tunedHasher() (*pwdhash.PasswordHasher, error) {
    argon2id := &argon2.Argon2idHasher{
        Memory:      128 * 1024,
        Iterations:  4,
        Parallelism: 2,
        SaltLength:  16,
        KeyLength:   32,
    }

    return pwdhash.New(pwdhash.WithHasher(argon2id))
}
```

You can also author alternate algorithms by satisfying the `pwdhash.Hasher` interface:

```go
type Hasher interface {
    ID() string
    Hash(password []byte) (string, error)
    Verify(password []byte, encoded string) (bool, error)
    NeedsRehash(encoded string) (bool, error)
}
```

Registering a hasher via `WithHasher` places it in the internal registry keyed by `ID()`. `Verify` and `NeedsRehash` parse the PHC string, look up the algorithm by its identifier, and dispatch to the correct implementation.

## PHC Encoding Basics

Internally the library parses/produces `encoding.EncodedHash` structures:

```
$argon2id$v=19$m=65536,t=3,p=4$<base64(salt)>$<base64(hash)>
```

- Algorithm identifiers map to registered hashers.
- Parameters are stored verbatim; `NeedsRehash` compares them to the current configuration to decide when to upgrade.

## Testing & Tooling

Run the full suite:

```bash
go test ./...
```

CI-style workflows use the Makefile helpers:

```bash
make lint   # golangci-lint run -v --fix
make test   # go test -covermode=count -coverprofile=count.out -v ./...
```

## Contributing

1. Fork and clone the repo.
2. Run `go test ./...` (and `make lint`) before submitting patches.
3. Keep exports documented, follow Go formatting (`gofmt`, `goimports`), and prefer table-driven tests.
4. Discuss larger API changes via issues or draft PRs to ensure alignment with the Argon2 focus of the project.

## License

This project is licensed under the MIT License. See `LICENSE` for details.
