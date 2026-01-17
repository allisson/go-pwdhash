package pwdhash

import (
	"testing"

	"github.com/allisson/go-pwdhash/argon2"
	"github.com/stretchr/testify/require"
)

type fakeHasher struct {
	id           string
	hash         string
	hashErr      error
	verifyResult bool
	verifyErr    error
	rehashNeeded bool
	rehashErr    error
}

func (f *fakeHasher) ID() string                          { return f.id }
func (f *fakeHasher) Hash([]byte) (string, error)         { return f.hash, f.hashErr }
func (f *fakeHasher) Verify([]byte, string) (bool, error) { return f.verifyResult, f.verifyErr }
func (f *fakeHasher) NeedsRehash(string) (bool, error)    { return f.rehashNeeded, f.rehashErr }

func TestPasswordHasher_HashUsesCurrentHasher(t *testing.T) {
	fake := &fakeHasher{id: "fake", hash: "encoded"}
	ph, err := New(WithHasher(fake))
	require.NoError(t, err)

	encoded, err := ph.Hash([]byte("pw"))
	require.NoError(t, err)
	require.Equal(t, "encoded", encoded)
}

func TestPasswordHasher_VerifyUnknownAlgorithm(t *testing.T) {
	ph, err := New(WithHasher(&fakeHasher{id: "fake"}))
	require.NoError(t, err)

	ok, err := ph.Verify([]byte("pw"), "$argon2id$v=19$m=1,t=1,p=1$YWJj$ZGVm")
	require.Error(t, err)
	require.False(t, ok)
}

func TestPasswordHasher_VerifyDelegatesToHasher(t *testing.T) {
	fake := &fakeHasher{
		id:           "argon2id",
		verifyResult: true,
	}
	ph, err := New(WithHasher(fake))
	require.NoError(t, err)

	ok, err := ph.Verify([]byte("pw"), "$argon2id$v=19$m=1,t=1,p=1$YWJj$ZGVm")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestPasswordHasher_NeedsRehashUnknownHasher(t *testing.T) {
	ph, err := New(WithHasher(&fakeHasher{id: "fake"}))
	require.NoError(t, err)

	needs, err := ph.NeedsRehash("$argon2id$v=19$m=1,t=1,p=1$YWJj$ZGVm")
	require.NoError(t, err)
	require.True(t, needs)
}

func TestPasswordHasher_NeedsRehashDelegates(t *testing.T) {
	fake := &fakeHasher{id: "argon2id", rehashNeeded: true}
	ph, err := New(WithHasher(fake))
	require.NoError(t, err)

	needs, err := ph.NeedsRehash("$argon2id$v=19$m=1,t=1,p=1$YWJj$ZGVm")
	require.NoError(t, err)
	require.True(t, needs)
}

func TestNewSetsDefaultHasher(t *testing.T) {
	ph, err := New()
	require.NoError(t, err)

	_, ok := ph.registry[argon2.Default().ID()]
	require.True(t, ok)
}
