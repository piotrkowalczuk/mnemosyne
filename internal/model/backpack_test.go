package model_test

import (
	"testing"

	"github.com/piotrkowalczuk/mnemosyne/internal/model"
)

func TestBackpack(t *testing.T) {
	b := model.Bag{}

	b.Set("key", "value")

	if !b.Has("key") {
		t.Errorf("bagpack should have specified key")
	}

	if b.Get("key") != "value" {
		t.Errorf("bagpack should have specified key")
	}
}

func TestBackpack_Scan(t *testing.T) {
	b := model.Bag{"A": "B"}
	val, err := b.Value()
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	b.Set("A", "C")
	if err = b.Scan(val); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if b.Get("A") != "B" {
		t.Errorf("wrong output, got %s", b.Get("A"))
	}

	exp := "unsupported data source type"
	b = model.Bag{}
	if err := b.Scan("data"); err != nil {
		if err.Error() != exp {
			t.Errorf("wrong error, expected %s but got %s", exp, err.Error())
		}
	} else {
		t.Error("error expected, got nil")
	}
}
