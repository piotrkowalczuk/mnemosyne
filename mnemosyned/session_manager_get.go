package mnemosyned

import (
	"time"

	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/jump"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type sessionManagerGet struct {
	storage storage
	cache   *cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (smg *sessionManagerGet) Get(ctx context.Context, req *mnemosynerpc.GetRequest) (*mnemosynerpc.GetResponse, error) {
	if req.AccessToken == "" {
		return nil, errMissingAccessToken
	}
	if node, ok := smg.cluster.GetOther(req.AccessToken); ok {
		smg.logger.Debug("get request forwarded", zap.String("to", node.Addr), zap.String("access_token", req.AccessToken))
		return node.Client.Get(ctx, req)
	}
	var (
		ses *mnemosynerpc.Session
		err error
	)

	hs := jump.Sum64(req.AccessToken)
	entry, ok := smg.cache.read(hs)
	if !ok || (!entry.refresh && time.Since(entry.exp) > smg.cache.ttl) {
		if ok {
			smg.cache.refresh(hs)
		}
		ses, err = smg.storage.Get(ctx, req.AccessToken)
		if err != nil {
			if err == errSessionNotFound && ok {
				smg.cache.del(hs)
			}
			return nil, err
		}
		smg.cache.put(hs, *ses)
	} else {
		ses = &entry.ses
	}

	return &mnemosynerpc.GetResponse{
		Session: ses,
	}, nil
}
