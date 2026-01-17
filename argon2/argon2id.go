package argon2

import (
	"crypto/rand"
	"fmt"
	"strconv"

	"golang.org/x/crypto/argon2"

	"github.com/allisson/go-pwdhash/internal/cast"
	"github.com/allisson/go-pwdhash/internal/encoding"
	"github.com/allisson/go-pwdhash/internal/subtle"
)

// Argon2idHasher wraps parameterized Argon2id hashing operations.
type Argon2idHasher struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// Default returns an Argon2idHasher configured with library defaults.
func Default() *Argon2idHasher {
	return &Argon2idHasher{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// ID reports the PHC algorithm identifier.
func (a *Argon2idHasher) ID() string {
	return "argon2id"
}

// Hash derives an Argon2id key and returns the PHC string.
func (a *Argon2idHasher) Hash(password []byte) (string, error) {
	salt := make([]byte, a.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := argon2.IDKey(
		password,
		salt,
		a.Iterations,
		a.Memory,
		a.Parallelism,
		a.KeyLength,
	)

	enc := encoding.EncodedHash{
		Algorithm: a.ID(),
		Version:   argon2.Version,
		Params: map[string]string{
			"m": strconv.Itoa(int(a.Memory)),
			"t": strconv.Itoa(int(a.Iterations)),
			"p": strconv.Itoa(int(a.Parallelism)),
		},
		Salt: salt,
		Hash: key,
	}

	return enc.String(), nil
}

// Verify recomputes the Argon2id hash and compares it in constant time.
func (a *Argon2idHasher) Verify(password []byte, encoded string) (bool, error) {
	parsed, err := encoding.Parse(encoded)
	if err != nil {
		return false, err
	}

	mem, err := cast.ConvertStringToUint32(parsed.Params["m"])
	if err != nil {
		return false, err
	}

	it, err := cast.ConvertStringToUint32(parsed.Params["t"])
	if err != nil {
		return false, err
	}

	par, err := cast.ConvertStringToUint8(parsed.Params["p"])
	if err != nil {
		return false, err
	}

	keyLen, err := cast.ConvertIntToUint32(len(parsed.Hash))
	if err != nil {
		return false, err
	}

	key := argon2.IDKey(
		password,
		parsed.Salt,
		it,
		mem,
		par,
		keyLen,
	)

	return subtle.ConstantTimeCompare(key, parsed.Hash), nil
}

// NeedsRehash reports whether the encoded parameters diverge from the current configuration.
func (a *Argon2idHasher) NeedsRehash(encoded string) (bool, error) {
	parsed, err := encoding.Parse(encoded)
	if err != nil {
		return false, err
	}

	if parsed.Algorithm != a.ID() {
		return true, nil
	}

	if parsed.Params["m"] != fmt.Sprint(a.Memory) {
		return true, nil
	}

	if parsed.Params["t"] != fmt.Sprint(a.Iterations) {
		return true, nil
	}

	return false, nil
}
