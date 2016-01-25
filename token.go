package mnemosyne

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/sha3"
	"golang.org/x/net/context"
)

// NewTokenContext returns a new Context that carries Token value.
func NewTokenContext(ctx context.Context, t Token) context.Context {
	return context.WithValue(ctx, TokenContextKey, t)
}

// TokenFromContext returns the Token value stored in context, if any.
func TokenFromContext(ctx context.Context) (Token, bool) {
	t, ok := ctx.Value(TokenContextKey).(Token)

	return t, ok
}

// Encode ...
func (t *Token) Encode() string {
	return string(t.Bytes())
}

// Bytes ...
func (t *Token) Bytes() []byte {
	if len(t.Key) < 10 {
		return t.Hash
	}

	return append(t.Key[:10], t.Hash...)
}

// DecodeToken parse string and allocates new token instance if ok. Expected token has format <key(10)><hash(n)>.
func DecodeToken(s []byte) (t Token) {
	if len(s) < 10 {
		return
	}

	return Token{
		Key:  bytes.TrimSpace(s[:10]),
		Hash: bytes.TrimSpace(s[10:]),
	}
}

// DecodeTokenString works like DecodeToken but accepts string.
func DecodeTokenString(s string) Token {
	return DecodeToken([]byte(s))
}

// NewToken ...
func NewToken(key, hash []byte) Token {
	if len(key) < 10 {
		return Token{
			Key:  append([]byte("0000000000")[:10-len(key)], key...),
			Hash: hash,
		}
	}
	return Token{
		Key:  key[:10],
		Hash: hash,
	}
}

// RandomToken ...
func RandomToken(generator RandomBytesGenerator, key []byte) (t Token, err error) {
	var buf []byte
	buf, err = generator.GenerateRandomBytes(128)
	if err != nil {
		return
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	hash := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(hash, buf)
	hash2 := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(hash2, hash)
	return NewToken(key, hash2), nil
}

// Value implements driver.Valuer interface.
func (t Token) Value() (driver.Value, error) {
	return string(t.Bytes()), nil
}

// Scan implements sql.Scanner interface.
func (t *Token) Scan(src interface{}) error {
	var (
		token Token
	)

	switch s := src.(type) {
	case []byte:
		token = DecodeToken(s)
	case string:
		token = DecodeTokenString(s)
	default:
		return errors.New("mnemosyne: token supports scan only from slice of bytes and string")
	}

	*t = token

	return nil
}

// IsEmpty ...
func (t *Token) IsEmpty() bool {
	if t == nil {
		return true
	}
	return len(t.Hash) == 0
}
