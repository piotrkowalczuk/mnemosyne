package mnemosyned

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
)

type sessionManagerDelete struct {
	spanner

	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (smd *sessionManagerDelete) Delete(ctx context.Context, req *mnemosynerpc.DeleteRequest) (*wrapperspb.Int64Value, error) {
	span, ctx := smd.span(ctx, "session-manager.delete")
	defer span.Finish()

	if req.AccessToken == "" && req.RefreshToken == "" && req.ExpireAtFrom == nil && req.ExpireAtTo == nil {
		return nil, status.Errorf(codes.InvalidArgument, "none of expected arguments was provided")
	}

	var expireAtFrom, expireAtTo *time.Time

	if req.ExpireAtFrom != nil {
		eaf := req.ExpireAtFrom.AsTime()
		expireAtFrom = &eaf
	}
	if req.ExpireAtTo != nil {
		eat := req.ExpireAtTo.AsTime()
		expireAtTo = &eat
	}

	aff, err := smd.storage.Delete(ctx, req.SubjectId, req.AccessToken, req.RefreshToken, expireAtFrom, expireAtTo)
	if err != nil {
		return nil, err
	}

	return &wrapperspb.Int64Value{Value: aff}, nil
}
