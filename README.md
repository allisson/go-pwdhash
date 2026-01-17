# pwdhash

`pwdhash` is a modern, extensible password hashing library for Go.

Inspired by:
- PHC string format
- argon2id best practices
- pwdlib (Python)

## Features

- Argon2id (default)
- PHC-compliant encoding
- Rehash detection
- Constant-time verification
- Extensible hasher registry

## Example

```go
hasher, _ := pwdhash.New()

hash, _ := hasher.Hash([]byte("secret"))

ok, _ := hasher.Verify([]byte("secret"), hash)

if ok {
    needs, _ := hasher.NeedsRehash(hash)
    if needs {
        hash, _ = hasher.Hash([]byte("secret"))
    }
}
