package mnemosynerpc

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/sha3"
	"golang.org/x/net/context"
)

// NewAccessTokenContext returns a new Context that carries token value.
func NewAccessTokenContext(ctx context.Context, at string) context.Context {
	return context.WithValue(ctx, accessTokenContextKey, at)
}

// AccessTokenFromContext returns the token value stored in context, if any.
func AccessTokenFromContext(ctx context.Context) (string, bool) {
	at, ok := ctx.Value(accessTokenContextKey).(string)

	return at, ok
}

// NewAccessToken allocates new access token based on given key and hash.
// Key should not be longer than 10 elements, otherwise will be truncated.
// If key is shorten then 10 elements, it will be filled with zeros at the beginning.
func NewAccessToken(key, hash string) string {
	if len(key) == 10 {
		return key + hash
	}
	if len(key) < 10 {
		return string([]byte("0000000000")[:10-len(key)]) + key + hash
	}
	return string(key[:10]) + hash
}

// RandomAccessToken generate Access Token with given key and generated hash of length 64.
func RandomAccessToken(key string) (at string, err error) {
	var buf []byte
	buf, err = generateRandomBytes(128)
	if err != nil {
		return
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	hash := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(hash, buf)
	hash2 := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(hash2, hash)
	return NewAccessToken(key, string(hash2)), nil
}

func generateRandomBytes(length int) ([]byte, error) {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil, err
	}
	return k, nil
}
