package mnemosynerpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	successTokens = map[string]AccessToken{
		"0000000key_hash": AccessToken{
			Key:  []byte("0000000key"),
			Hash: []byte("_hash"),
		},
		"0000000001be29318725cdbd51aa0078dee7ff5b4eda7f18d2d91c87b1117f2f2d35db044eceaf19dfbd9c56ec14fd5d7aa8e795e3c25e288eee0bef7b7828fbe5710f2362": AccessToken{
			Key:  []byte("0000000001"),
			Hash: []byte("be29318725cdbd51aa0078dee7ff5b4eda7f18d2d91c87b1117f2f2d35db044eceaf19dfbd9c56ec14fd5d7aa8e795e3c25e288eee0bef7b7828fbe5710f2362"),
		},
	}
)

func TestAccessToken_Scan(t *testing.T) {
	for given, expected := range successTokens {
		token := &AccessToken{}

		err := token.Scan(given)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, *token)
		}
	}
}

func TestAccessToken_Value(t *testing.T) {
	for expected, given := range successTokens {
		v, err := given.Value()
		if assert.NoError(t, err) {
			assert.Equal(t, expected, v.(string))
		}
	}
}

func TestAccessToken_Bytes(t *testing.T) {
	for expected, given := range successTokens {
		v := given.Bytes()
		assert.Equal(t, expected, string(v))
	}
}

func TestDecodeTokenString(t *testing.T) {
	success := map[string]AccessToken{
		"01234567895234532534523": AccessToken{Key: []byte("0123456789"), Hash: []byte("5234532534523")},
		"01234567891":             AccessToken{Key: []byte("0123456789"), Hash: []byte("1")},
		"0123456789a":             AccessToken{Key: []byte("0123456789"), Hash: []byte("a")},
		"1":                       AccessToken{Key: nil, Hash: nil},
		":1":                      AccessToken{Key: nil, Hash: nil},
		"1:":                      AccessToken{Key: nil, Hash: nil},
		":":                       AccessToken{Key: nil, Hash: nil},
		"":                        AccessToken{Key: nil, Hash: nil},
		":   ":                    AccessToken{Key: nil, Hash: nil},
		"   :":                    AccessToken{Key: nil, Hash: nil},
	}

	for given, expected := range success {
		got := DecodeAccessTokenString(given)
		assert.Equal(t, expected, got)
	}
}

func TestRandomToken(t *testing.T) {
	token, err := RandomAccessToken([]byte("abc"))

	if assert.NoError(t, err) {
		assert.Len(t, token.Hash, 128)
		assert.Len(t, token.Key, 10)
		assert.Equal(t, token.Key, []byte("0000000abc"))
	}
}
