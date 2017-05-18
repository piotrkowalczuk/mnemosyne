package mnemosyned

import (
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type sessionManagerExists struct {
	storage storage
	cache   *cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (sme *sessionManagerExists) Exists(ctx context.Context, req *mnemosynerpc.ExistsRequest) (*mnemosynerpc.ExistsResponse, error) {
	if req.AccessToken == "" {
		return nil, errMissingAccessToken
	}
	if node, ok := sme.cluster.GetOther(req.AccessToken); ok {
		sme.logger.Debug("exists request forwarded", zap.String("to", node.Addr), zap.String("access_token", req.AccessToken))
		return node.Client.Exists(ctx, req)
	}

	exists, err := sme.storage.Exists(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	return &mnemosynerpc.ExistsResponse{
		Exists: exists,
	}, nil
}
