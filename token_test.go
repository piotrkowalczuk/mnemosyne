package mnemosyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	successTokens = map[string]Token{
		"0000000key_hash": Token{
			Key:  []byte("0000000key"),
			Hash: []byte("_hash"),
		},
		"0000000001be29318725cdbd51aa0078dee7ff5b4eda7f18d2d91c87b1117f2f2d35db044eceaf19dfbd9c56ec14fd5d7aa8e795e3c25e288eee0bef7b7828fbe5710f2362": Token{
			Key:  []byte("0000000001"),
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
			assert.Equal(t, expected, string(v.([]uint8)))
		}
	}
}

func TestDecodeTokenString(t *testing.T) {
	success := map[string]Token{
		"01234567895234532534523": Token{Key: []byte("0123456789"), Hash: []byte("5234532534523")},
		"01234567891":             Token{Key: []byte("0123456789"), Hash: []byte("1")},
		"0123456789a":             Token{Key: []byte("0123456789"), Hash: []byte("a")},
		"1":                       Token{Key: nil, Hash: nil},
		":1":                      Token{Key: nil, Hash: nil},
		"1:":                      Token{Key: nil, Hash: nil},
		":":                       Token{Key: nil, Hash: nil},
		"":                        Token{Key: nil, Hash: nil},
		":   ":                    Token{Key: nil, Hash: nil},
		"   :":                    Token{Key: nil, Hash: nil},
	}

	for given, expected := range success {
		got := DecodeTokenString(given)
		assert.Equal(t, expected, got)
	}
}

func TestRandomToken(t *testing.T) {
	token, err := RandomToken(&SystemRandomBytesGenerator{}, []byte("abc"))

	if assert.NoError(t, err) {
		assert.Len(t, token.Hash, 64)
		assert.Len(t, token.Key, 10)
		assert.Equal(t, token.Key, []byte("0000000abc"))
	}
}
