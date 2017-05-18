package mnemosyne

import "testing"

func TestRandomToken(t *testing.T) {
	token, err := RandomAccessToken()

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
		at, err := RandomAccessToken()
		if err != nil {
			b.Fatalf("unexpected error: %s", err.Error())
		}
		benchAccessToken = at
	}
}
