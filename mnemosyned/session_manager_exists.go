package mnemosyned

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type sessionManagerExists struct {
	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (sme *sessionManagerExists) Exists(ctx context.Context, req *mnemosynerpc.ExistsRequest) (*wrappers.BoolValue, error) {
	if req.AccessToken == "" {
		return nil, errMissingAccessToken
	}
	if node, ok := sme.cluster.GetOther(req.AccessToken); ok {
		sme.logger.Debug("exists request forwarded", zap.String("remote_addr", node.Addr), zap.String("access_token", req.AccessToken))
		return node.Client.Exists(ctx, req)
	}

	exists, err := sme.storage.Exists(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	return &wrappers.BoolValue{Value: exists}, nil
}
