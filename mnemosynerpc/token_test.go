package mnemosynerpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomToken(t *testing.T) {
	token, err := RandomAccessToken("abc")

	if assert.NoError(t, err) {
		assert.Len(t, string(token), 138)
		assert.Equal(t, string(token[:10]), "0000000abc")
	}
}
