package mnemosyned

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
)

type sessionManagerList struct {
	storage storage
}

func (sml *sessionManagerList) List(ctx context.Context, req *mnemosynerpc.ListRequest) (*mnemosynerpc.ListResponse, error) {
	var (
		expireAtFrom, expireAtTo *time.Time
	)
	if req.ExpireAtFrom != nil {
		eaf, err := ptypes.Timestamp(req.ExpireAtFrom)
		if err != nil {
			return nil, err
		}
		expireAtFrom = &eaf
	}
	if req.ExpireAtTo != nil {
		eat, err := ptypes.Timestamp(req.ExpireAtTo)
		if err != nil {
			return nil, err
		}
		expireAtTo = &eat
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	sessions, err := sml.storage.List(ctx, req.Offset, req.Limit, expireAtFrom, expireAtTo)
	if err != nil {
		return nil, err
	}
	return &mnemosynerpc.ListResponse{
		Sessions: sessions,
	}, nil
}
