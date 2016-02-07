package mnemosynetest_test

import (
	"testing"

	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
)

func TestMnemosyne(t *testing.T) {
	var mock interface{} = &mnemosynetest.Mnemosyne{}

	if _, ok := mock.(mnemosyne.Mnemosyne); !ok {
		t.Errorf("mock should implement original interface, but does not")
	}
}
