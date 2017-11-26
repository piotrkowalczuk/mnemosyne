package mnemosyned

import (
	"time"

	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/jump"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type sessionManagerGet struct {
	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (smg *sessionManagerGet) Get(ctx context.Context, req *mnemosynerpc.GetRequest) (*mnemosynerpc.GetResponse, error) {
	if req.AccessToken == "" {
		return nil, errMissingAccessToken
	}
	if node, ok := smg.cluster.GetOther(req.AccessToken); ok {
		smg.logger.Debug("get request forwarded", zap.String("remote_addr", node.Addr), zap.String("access_token", req.AccessToken))
		return node.Client.Get(ctx, req)
	}
	var (
		ses *mnemosynerpc.Session
		err error
	)

	hs := jump.Sum64(req.AccessToken)
	entry, ok := smg.cache.Read(hs)
	if !ok || (!entry.Refresh && time.Since(entry.Exp) > smg.cache.TTL) {
		if ok {
			smg.cache.Refresh(hs)
		}
		ses, err = smg.storage.Get(ctx, req.AccessToken)
		if err != nil {
			if err == storage.ErrSessionNotFound && ok {
				smg.cache.Del(hs)
			}
			return nil, err
		}
		smg.cache.Put(hs, *ses)
	} else {
		ses = &entry.Ses
	}

	return &mnemosynerpc.GetResponse{
		Session: ses,
	}, nil
}
