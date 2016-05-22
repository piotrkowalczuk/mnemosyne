package mnemosynerpc

import (
	"net/url"
	"time"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/oauth2"
)

const (
	// TokenContextKey ...
	TokenContextKey = "mnemosyne_token"
	// AccessTokenContextKey is used by Mnemosyne internally to retrieve session token from context.Context.
	AccessTokenContextKey = "mnemosyne_access_token"
	// AccessTokenMetadataKey is used by Mnemosyne to retrieve session token from gRPC metadata object.
	AccessTokenMetadataKey = "authorization"
)

// Token implements oauth2.TokenSource interface.
func (s *Session) Token() (*oauth2.Token, error) {
	var (
		err      error
		expireAt time.Time
	)
	if s.ExpireAt != nil {
		expireAt, err = ptypes.Timestamp(s.ExpireAt)
		if err != nil {
			return nil, err
		}
	}
	token := &oauth2.Token{
		AccessToken: s.AccessToken.Encode(),
		Expiry:      expireAt,
	}
	if s.Bag != nil && len(s.Bag) > 0 {
		token = token.WithExtra(bagToURLValues(s.Bag))
	}

	return token, nil
}

func bagToURLValues(b map[string]string) url.Values {
	r := make(map[string][]string, len(b))
	for k, v := range b {
		r[k] = []string{v}
	}
	return url.Values(r)
}
