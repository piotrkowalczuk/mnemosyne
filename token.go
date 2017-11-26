package mnemosyne

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/sha3"
	"golang.org/x/net/context"
)

const (
	// AccessTokenMetadataKey is used by Mnemosyne to retrieve session token from gRPC metadata object.
	AccessTokenMetadataKey = "authorization"
)

type key struct{}

var accessTokenContextKey = key{}

// NewAccessTokenContext returns a new Context that carries token value.
func NewAccessTokenContext(ctx context.Context, at string) context.Context {
	return context.WithValue(ctx, accessTokenContextKey, at)
}

// AccessTokenFromContext returns the token value stored in context, if any.
func AccessTokenFromContext(ctx context.Context) (string, bool) {
	at, ok := ctx.Value(accessTokenContextKey).(string)

	return at, ok
}

// RandomAccessToken generate Access Token with given key and generated hash of length 64.
func RandomAccessToken() (string, error) {
	buf, err := generateRandomBytes(128)
	if err != nil {
		return "", err
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	hash := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(hash, buf)
	hash2 := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(hash2, hash)
	return string(hash2), nil
}

func generateRandomBytes(length int) ([]byte, error) {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil, err
	}
	return k, nil
}
