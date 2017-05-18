package mnemosyned

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type sessionManagerDelete struct {
	storage storage
	cache   *cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (smd *sessionManagerDelete) Delete(ctx context.Context, req *mnemosynerpc.DeleteRequest) (*mnemosynerpc.DeleteResponse, error) {
	if req.AccessToken == "" && req.RefreshToken == "" && req.ExpireAtFrom == nil && req.ExpireAtTo == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "none of expected arguments was provided")
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

	affected, err := smd.storage.Delete(ctx, req.SubjectId, req.AccessToken, req.RefreshToken, expireAtFrom, expireAtTo)
	if err != nil {
		return nil, err
	}

	return &mnemosynerpc.DeleteResponse{
		Count: affected,
	}, nil
}
