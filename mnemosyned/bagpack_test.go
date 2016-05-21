package mnemosyned

import "testing"

func TestBagpack(t *testing.T) {
	b := bagpack{}

	b.Set("key", "value")

	if !b.Has("key") {
		t.Errorf("bagpack should have specified key")
	}

	if b.Get("key") != "value" {
		t.Errorf("bagpack should have specified key")
	}
}
