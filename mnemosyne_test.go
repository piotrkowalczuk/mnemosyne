package mnemosyne

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/oauth2"
)

func testTimestampProtoDate(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) *tspb.Timestamp {
	t, _ := ptypes.TimestampProto(time.Date(year, month, day, hour, min, sec, nsec, loc))

	return t
}

func TestSession_Token(t *testing.T) {
	cases := []struct {
		given    *Session
		expected oauth2.Token
	}{
		{
			given: &Session{
				AccessToken: &AccessToken{Key: []byte("0000000001"), Hash: []byte("abc")},
				ExpireAt:    testTimestampProtoDate(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: oauth2.Token{
				AccessToken: "0000000001abc",
				Expiry:      time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			given: &Session{
				AccessToken: &AccessToken{Key: []byte("0000000001"), Hash: []byte("999999999999")},
				Bag: map[string]string{
					"firstName": "John",
					"lastName":  "Snow",
				},
			},
			expected: oauth2.Token{
				AccessToken: "0000000001999999999999",
			},
		},
	}

	for i, c := range cases {
		got, err := c.given.Token()
		if err != nil {
			t.Errorf("unexpected error for #%d: %s", i, err.Error())
			continue
		}
		if c.expected.AccessToken != got.AccessToken {
			t.Errorf("wrong access token, expected %s but got %s", c.expected.AccessToken, got.AccessToken)
		}
		if c.expected.Type() != got.Type() {
			t.Errorf("wrong token type, expected %s but got %s", c.expected.Type(), got.Type())
		}
		if c.expected.RefreshToken != got.RefreshToken {
			t.Errorf("wrong refresh token, expected %s but got %s", c.expected.RefreshToken, got.RefreshToken)
		}
		if c.expected.Expiry != got.Expiry {
			t.Errorf("wrong expiry, expected %s but got %s", c.expected.Expiry, got.Expiry)
		}
		for k, v := range c.given.Bag {
			if s, ok := got.Extra(k).(string); !ok || s != v {
				t.Errorf("wrong bag value, expected %s but got %s for key: %s", v, s, k)
			}
		}
	}
}
