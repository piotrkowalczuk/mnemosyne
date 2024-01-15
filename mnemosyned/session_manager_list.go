package mnemosyned

import (
	"time"

	"golang.org/x/net/context"

	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
)

type sessionManagerList struct {
	spanner

	storage storage.Storage
}

func (sml *sessionManagerList) List(ctx context.Context, req *mnemosynerpc.ListRequest) (*mnemosynerpc.ListResponse, error) {
	span, ctx := sml.span(ctx, "session-manager.list")
	defer span.Finish()

	var expireAtFrom, expireAtTo *time.Time
	if req.GetQuery().GetExpireAtFrom() != nil {
		eaf := req.GetQuery().GetExpireAtFrom().AsTime()
		expireAtFrom = &eaf
	}
	if req.GetQuery().GetExpireAtTo() != nil {
		eat := req.GetQuery().GetExpireAtTo().AsTime()
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
