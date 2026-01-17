package pwdhash

type Hasher interface {
	ID() string
	Hash(password []byte) (string, error)
	Verify(password []byte, encoded string) (bool, error)
	NeedsRehash(encoded string) (bool, error)
}
