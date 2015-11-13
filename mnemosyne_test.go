package mnemosyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	successTokens = map[string]Token{
		"1:1": Token{
			Key:  "1",
			Hash: "1",
		},
		"1:be29318725cdbd51aa0078dee7ff5b4eda7f18d2d91c87b1117f2f2d35db044eceaf19dfbd9c56ec14fd5d7aa8e795e3c25e288eee0bef7b7828fbe5710f2362": Token{
			Key:  "1",
			Hash: "be29318725cdbd51aa0078dee7ff5b4eda7f18d2d91c87b1117f2f2d35db044eceaf19dfbd9c56ec14fd5d7aa8e795e3c25e288eee0bef7b7828fbe5710f2362",
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
			assert.Equal(t, expected, v)
		}
	}
}
