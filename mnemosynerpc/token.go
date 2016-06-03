package mnemosynerpc

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/sha3"
	"golang.org/x/net/context"
)

// NewAccessTokenContext returns a new Context that carries Token value.
func NewAccessTokenContext(ctx context.Context, at AccessToken) context.Context {
	return context.WithValue(ctx, accessTokenContextKey, at)
}

// AccessTokenFromContext returns the Token value stored in context, if any.
func AccessTokenFromContext(ctx context.Context) (AccessToken, bool) {
	at, ok := ctx.Value(accessTokenContextKey).(AccessToken)

	return at, ok
}

// Encode ...
func (at AccessToken) Encode() string {
	return string(at.Bytes())
}

// Bytes ...
func (at AccessToken) Bytes() []byte {
	if len(at.Key) < 10 {
		return at.Hash
	}

	return append(at.Key[:10], at.Hash...)
}

// DecodeAccessToken parse byte slice and allocates new token instance if ok.
// Expected token has format <key(10)><hash(n)>.
func DecodeAccessToken(s []byte) (at AccessToken) {
	if len(s) <= 10 {
		return
	}

	return AccessToken{
		Key:  bytes.TrimSpace(s[:10]),
		Hash: bytes.TrimSpace(s[10:]),
	}
}

// ParseAccessToken parse given string similar way how DecodeAccessToken does with byte slice.
func ParseAccessToken(s string) AccessToken {
	return DecodeAccessToken([]byte(s))
}

// DecodeAccessTokenString works like DecodeToken but accepts string.
func DecodeAccessTokenString(s string) AccessToken {
	return DecodeAccessToken([]byte(s))
}

// NewAccessToken allocates new access token based on given key and hash.
// Key should not be longer than 10 elements, otherwise will be truncated.
// If key is shorten then 10 elements, it will be filled with zeros at the beginning.
func NewAccessToken(key, hash []byte) AccessToken {
	if len(key) < 10 {
		return AccessToken{
			Key:  append([]byte("0000000000")[:10-len(key)], key...),
			Hash: hash,
		}
	}
	return AccessToken{
		Key:  key[:10],
		Hash: hash,
	}
}

// Value implements driver.Valuer interface.
func (at AccessToken) Value() (driver.Value, error) {
	return string(at.Bytes()), nil
}

// Scan implements sql.Scanner interface.
func (at *AccessToken) Scan(src interface{}) error {
	var (
		token AccessToken
	)

	switch s := src.(type) {
	case []byte:
		token = DecodeAccessToken(s)
	case string:
		token = DecodeAccessTokenString(s)
	default:
		return errors.New("mnemosyne: token supports scan only from slice of bytes and string")
	}

	*at = token

	return nil
}

// IsEmpty returns true if hash is zero length.
func (at AccessToken) IsEmpty() bool {
	return len(at.Hash) == 0
}

// RandomAccessToken generate Access Token with given key and generated hash of length 64.
func RandomAccessToken(key []byte) (at AccessToken, err error) {
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
	return NewAccessToken(key, hash2), nil
}

func generateRandomBytes(length int) ([]byte, error) {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil, err
	}
	return k, nil
}
