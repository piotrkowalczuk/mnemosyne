package mnemosyned

import "testing"

func TestSystemRandomBytesGenerator_generateRandomBytes(t *testing.T) {
	g := systemRandomBytesGenerator{}
	for i := 1; i < 11; i++ {
		got, err := g.generateRandomBytes(i)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		if len(got) != i {
			t.Fatalf("wrong output length, expected %d but got %d", i, len(got))
		}
	}
}
