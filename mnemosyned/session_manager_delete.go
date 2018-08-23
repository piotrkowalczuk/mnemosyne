package mnemosyned

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type sessionManagerDelete struct {
	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (smd *sessionManagerDelete) Delete(ctx context.Context, req *mnemosynerpc.DeleteRequest) (*wrappers.Int64Value, error) {
	if req.AccessToken == "" && req.RefreshToken == "" && req.ExpireAtFrom == nil && req.ExpireAtTo == nil {
		return nil, status.Errorf(codes.InvalidArgument, "none of expected arguments was provided")
	}

	var expireAtFrom, expireAtTo *time.Time

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

	aff, err := smd.storage.Delete(ctx, req.SubjectId, req.AccessToken, req.RefreshToken, expireAtFrom, expireAtTo)
	if err != nil {
		return nil, err
	}

	return &wrappers.Int64Value{Value: aff}, nil
}
