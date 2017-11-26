package mnemosyne_test

import (
	"context"
	"testing"

	"github.com/piotrkowalczuk/mnemosyne"
)

func TestNewAccessTokenContext(t *testing.T) {
	exp := "123"
	ctx := mnemosyne.NewAccessTokenContext(context.Background(), exp)
	if got, ok := mnemosyne.AccessTokenFromContext(ctx); ok {
		if exp != got {
			t.Errorf("wrong access token, expected %s but got %s", exp, got)
		}
		return
	}
	t.Error("missing access token")
}

func TestRandomToken(t *testing.T) {
	token, err := mnemosyne.RandomAccessToken()

	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if len(token) != 128 {
		t.Errorf("wrong length, expected %d but got %d", 128, len(token))
	}
}

var (
	benchAccessToken string
)

func BenchmarkRandomAccessToken(b *testing.B) {
	bn := int32(b.N)

	b.ResetTimer()
	for n := int32(1); n < bn; n++ {
		at, err := mnemosyne.RandomAccessToken()
		if err != nil {
			b.Fatalf("unexpected error: %s", err.Error())
		}
		benchAccessToken = at
	}
}
