package mnemosyne

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	successTokens = map[string]Token{
		hex.EncodeToString([]byte("1:1")): Token{
			Key:  []byte("1"),
			Hash: []byte("1"),
		},
		hex.EncodeToString([]byte("1:be29318725cdbd51aa0078dee7ff5b4eda7f18d2d91c87b1117f2f2d35db044eceaf19dfbd9c56ec14fd5d7aa8e795e3c25e288eee0bef7b7828fbe5710f2362")): Token{
			Key:  []byte("1"),
			Hash: []byte("be29318725cdbd51aa0078dee7ff5b4eda7f18d2d91c87b1117f2f2d35db044eceaf19dfbd9c56ec14fd5d7aa8e795e3c25e288eee0bef7b7828fbe5710f2362"),
		},
	}
)

func TestToken_Scan(t *testing.T) {
	for given, expected := range successTokens {
		token := &Token{}

		err := token.Scan(given)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, *token)
		}
	}
}

func TestToken_Value(t *testing.T) {
	for expected, given := range successTokens {
		v, err := given.Value()
		if assert.NoError(t, err) {
			h, _ := hex.DecodeString(expected)
			assert.Equal(t, h, v)
		}
	}
}

func TestDecodeTokenString(t *testing.T) {
	success := map[string]Token{
		hex.EncodeToString([]byte("1:5234532534523")): Token{Key: []byte("1"), Hash: []byte("5234532534523")},
		hex.EncodeToString([]byte("1:1")):             Token{Key: []byte("1"), Hash: []byte("1")},
		hex.EncodeToString([]byte("1")):               Token{Key: nil, Hash: []byte("1")},
		hex.EncodeToString([]byte(":1")):              Token{Key: nil, Hash: []byte("1")},
		hex.EncodeToString([]byte("1:")):              Token{Key: nil, Hash: []byte("1")},
		hex.EncodeToString([]byte(":")):               Token{Key: nil, Hash: nil},
		hex.EncodeToString([]byte("")):                Token{Key: nil, Hash: nil},
		hex.EncodeToString([]byte(":   ")):            Token{Key: nil, Hash: nil},
		hex.EncodeToString([]byte("   :")):            Token{Key: nil, Hash: nil},
	}

	for given, expected := range success {
		got, err := DecodeTokenString(given)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, got)
		}
	}
}

func TestRandomToken(t *testing.T) {
	token, err := RandomToken(&SystemRandomBytesGenerator{}, []byte(hex.EncodeToString([]byte("1"))))

	if assert.NoError(t, err) {
		assert.Len(t, token.Hash, 128)
		assert.Len(t, token.Key, 4)
	}
}
