package mnemosyned

import (
	"crypto/rand"
	"io"
)

type randomBytesGenerator interface {
	generateRandomBytes(int) ([]byte, error)
}

type systemRandomBytesGenerator struct {
}

// GenerateRandomBytes creates a random key with the given length in bytes.
func (srbg *systemRandomBytesGenerator) generateRandomBytes(length int) ([]byte, error) {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil, err
	}
	return k, nil
}
