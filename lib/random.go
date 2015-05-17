package lib

import (
	"crypto/rand"
	"io"
)

// RandomBytesGenerator ...
type RandomBytesGenerator interface {
	GenerateRandomBytes(int) []byte
}

// SystemRandomBytesGenerator ...
type SystemRandomBytesGenerator struct {
}

// SystemRandomBytesGenerator creates a random key with the given length in bytes.
func (srbg *SystemRandomBytesGenerator) GenerateRandomBytes(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
