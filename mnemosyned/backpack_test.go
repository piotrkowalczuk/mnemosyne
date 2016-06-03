package mnemosyned

import "testing"

func TestBackpack(t *testing.T) {
	b := bag{}

	b.set("key", "value")

	if !b.has("key") {
		t.Errorf("bagpack should have specified key")
	}

	if b.get("key") != "value" {
		t.Errorf("bagpack should have specified key")
	}
}
