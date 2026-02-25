package pwdhash_test

import (
	"fmt"
	"strings"

	"github.com/allisson/go-pwdhash"
)

func ExampleNew() {
	hasher, err := pwdhash.New()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Default hasher created: %T\n", hasher)
	// Output: Default hasher created: *pwdhash.PasswordHasher
}

func ExamplePasswordHasher_Hash() {
	hasher, err := pwdhash.New(pwdhash.WithPolicy(pwdhash.PolicyInteractive))
	if err != nil {
		panic(err)
	}

	// Hash a password
	encoded, err := hasher.Hash([]byte("s3cret"))
	if err != nil {
		panic(err)
	}

	// PHC hashes for Argon2id start with $argon2id$
	fmt.Println(strings.HasPrefix(encoded, "$argon2id$"))
	// Output: true
}

func ExamplePasswordHasher_Verify() {
	hasher, err := pwdhash.New(pwdhash.WithPolicy(pwdhash.PolicyInteractive))
	if err != nil {
		panic(err)
	}

	password := []byte("s3cret")

	// Generate a hash to verify
	encoded, err := hasher.Hash(password)
	if err != nil {
		panic(err)
	}

	// Verify the correct password
	ok, err := hasher.Verify([]byte("s3cret"), encoded)
	if err != nil {
		panic(err)
	}
	fmt.Println("Correct password:", ok)

	// Verify an incorrect password
	ok, err = hasher.Verify([]byte("wrong_password"), encoded)
	if err != nil {
		panic(err)
	}
	fmt.Println("Incorrect password:", ok)
	// Output:
	// Correct password: true
	// Incorrect password: false
}

func ExamplePasswordHasher_NeedsRehash() {
	// 1. Simulate an old hash created with the Interactive policy
	oldHasher, err := pwdhash.New(pwdhash.WithPolicy(pwdhash.PolicyInteractive))
	if err != nil {
		panic(err)
	}
	encoded, err := oldHasher.Hash([]byte("s3cret"))
	if err != nil {
		panic(err)
	}

	// 2. Initialize a new hasher with a stronger Moderate policy
	newHasher, err := pwdhash.New(pwdhash.WithPolicy(pwdhash.PolicyModerate))
	if err != nil {
		panic(err)
	}

	// 3. Check if the old hash needs to be upgraded using the new hasher
	needsRehash, err := newHasher.NeedsRehash(encoded)
	if err != nil {
		panic(err)
	}

	fmt.Println("Needs rehash with stronger policy:", needsRehash)
	// Output: Needs rehash with stronger policy: true
}

func ExampleWithPolicy() {
	// Instantiate a hasher with a specific policy
	hasher, err := pwdhash.New(
		pwdhash.WithPolicy(pwdhash.PolicyModerate),
	)
	if err != nil {
		panic(err)
	}

	// Use standard operations
	encoded, err := hasher.Hash([]byte("my_secure_password"))
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.HasPrefix(encoded, "$argon2id$"))
	// Output: true
}
