package mnemosyned

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
)

type sessionManagerList struct {
	spanner

	storage storage.Storage
}

func (sml *sessionManagerList) List(ctx context.Context, req *mnemosynerpc.ListRequest) (*mnemosynerpc.ListResponse, error) {
	span, ctx := sml.span(ctx, "session-manager.list")
	defer span.Finish()

	var (
		expireAtFrom, expireAtTo *time.Time
	)
	if req.GetQuery().GetExpireAtFrom() != nil {
		eaf, err := ptypes.Timestamp(req.GetQuery().GetExpireAtFrom())
		if err != nil {
			return nil, err
		}
		expireAtFrom = &eaf
	}
	if req.GetQuery().GetExpireAtTo() != nil {
		eat, err := ptypes.Timestamp(req.GetQuery().GetExpireAtTo())
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
