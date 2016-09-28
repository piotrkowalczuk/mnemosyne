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

//func TestNewAccessToken(t *testing.T) {
//	cases := map[string]struct {
//		expected string
//		key      int32
//		hash     string
//	}{
//		"one": {
//			expected: "0000000001f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//			key:      1,
//			hash:     "f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//		},
//		"two": {
//			expected: "0000000031f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//			key:      31,
//			hash:     "f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//		},
//		"tree": {
//			expected: "0000000123f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//			key:      123,
//			hash:     "f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//		},
//		"four": {
//			expected: "0000001234f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//			key:      1234,
//			hash:     "f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//		},
//	}
//
//	for hint, c := range cases {
//		t.Run(hint, func(t *testing.T) {
//			at := NewAccessToken(c.key, c.hash)
//			if len(at) != 138 {
//				t.Fatalf("wrong access token length, expected %d but got %d", 138, len(at))
//			}
//			key, hash, valid := SplitAccessToken(at)
//			if !valid {
//				t.Fatal("access token split failure")
//			}
//
//			if c.key != key {
//				t.Errorf("invalid key, expected %d but got %d", c.key, key)
//			}
//			if c.hash != hash {
//				t.Errorf("invalid hash, expected %d but got %d", c.hash, hash)
//			}
//
//		})
//	}
//}

//func TestSplitAccessToken(t *testing.T) {
//	cases := map[string]struct {
//		given string
//		key   int32
//		hash  string
//		valid bool
//	}{
//		"basic": {
//			given: "0000000031f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//			key:   31,
//			hash:  "f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414",
//			valid: true,
//		},
//	}
//
//	for hint, c := range cases {
//		t.Run(hint, func(t *testing.T) {
//			key, hash, valid := SplitAccessToken(c.given)
//			if c.valid != valid {
//				t.Fatalf("invalid validity, expected %t but got %t", c.valid, valid)
//			}
//			if c.valid {
//				if c.key != key {
//					t.Errorf("invalid key, expected %d but got %d", c.key, key)
//				}
//				if c.hash != hash {
//					t.Errorf("invalid hash, expected %d but got %d", c.hash, hash)
//				}
//			}
//		})
//	}
//}

var (
	benchAccessToken string
	benchKey         int32
	benchHash        string
	benchValid       bool
)

//
//func BenchmarkNewAccessToken(b *testing.B) {
//	bn := int32(b.N)
//	h := "1234567890"
//
//	b.ResetTimer()
//	for n := int32(0); n < bn; n++ {
//		at := NewAccessToken(n, h)
//		benchAccessToken = at
//	}
//}

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

//func BenchmarkSplitAccessToken(b *testing.B) {
//	s := "0000000031f18c69d4799e4617563420537e6fc4ce58eb08f9fa74b33dc043d1dc3e10c8ceb709cf6d84a304d614384e6d11d7786670f17511ce55f4e6ec60bfd0c64c7414"
//	b.ResetTimer()
//	for n := 0; n < b.N; n++ {
//		key, hash, valid := SplitAccessToken(s)
//		benchAccessToken = hash
//		benchKey = key
//		benchValid = valid
//	}
//}
