package mnemosyned

import (
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

type sessionManagerExists struct {
	spanner

	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (sme *sessionManagerExists) Exists(ctx context.Context, req *mnemosynerpc.ExistsRequest) (*wrappers.BoolValue, error) {
	span, ctx := sme.span(ctx, "session-manager.exists")
	defer span.Finish()

	if req.AccessToken == "" {
		return nil, errMissingAccessToken
	}
	if node, ok := sme.cluster.GetOther(req.AccessToken); ok {
		if cluster.IsInternalRequest(ctx) {
			return nil, status.Errorf(codes.FailedPrecondition,
				"it should be final destination of exists request (%s), but found another node for it: %s",
				req.GetAccessToken(),
				node.Addr,
			)
		}
		sme.logger.Debug("exists request forwarded", zap.String("remote_addr", node.Addr), zap.String("access_token", req.AccessToken))
		return node.Client.Exists(ctx, req)
	}

	exists, err := sme.storage.Exists(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	return &wrappers.BoolValue{Value: exists}, nil
}
