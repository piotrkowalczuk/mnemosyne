package mnemosyned

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/jump"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	// DefaultTTL is session time to live default value.
	DefaultTTL = 24 * time.Minute
	// DefaultTTC is time to clear default value.
	DefaultTTC = 1 * time.Minute
	// DefaultCacheSize determines how big cache should be at the beginning.
	DefaultCacheSize = 100000
)

var (
	errSessionNotFound    = grpc.Errorf(codes.NotFound, "mnemosyned: session not found")
	errMissingAccessToken = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing access token")
	errMissingSubjectID   = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing subject id")
	errMissingSession     = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing session")
)

type sessionManagerOpts struct {
	addr       string
	cluster    *cluster.Cluster
	ttc        time.Duration
	logger     log.Logger
	storage    storage
	monitoring *monitoring
}

type cacheEntry struct {
	ses     mnemosynerpc.Session
	exp     time.Time
	refresh bool
}

type sessionManager struct {
	addr    string
	ttc     time.Duration
	logger  log.Logger
	monitor *monitoring
	storage storage
	cluster *cluster.Cluster
	alloc   struct {
		abandon  handlerFunc
		context  handlerFunc
		delete   handlerFunc
		exists   handlerFunc
		get      handlerFunc
		list     handlerFunc
		setValue handlerFunc
		start    handlerFunc
	}
	// cache
	ttl      time.Duration
	data     map[uint64]*cacheEntry
	dataLock sync.RWMutex
}

func newSessionManager(opts sessionManagerOpts) (*sessionManager, error) {
	return &sessionManager{
		addr:    opts.addr,
		cluster: opts.cluster,
		ttc:     opts.ttc,
		alloc: struct {
			abandon  handlerFunc
			context  handlerFunc
			delete   handlerFunc
			exists   handlerFunc
			get      handlerFunc
			list     handlerFunc
			setValue handlerFunc
			start    handlerFunc
		}{
			abandon:  newHandlerFunc("abandon"),
			context:  newHandlerFunc("context"),
			delete:   newHandlerFunc("delete"),
			exists:   newHandlerFunc("exists"),
			get:      newHandlerFunc("get"),
			list:     newHandlerFunc("list"),
			setValue: newHandlerFunc("set_value"),
			start:    newHandlerFunc("start"),
		},
		logger:  opts.logger,
		storage: opts.storage,
		monitor: opts.monitoring,
		ttl:     5 * time.Second,
		data:    make(map[uint64]*cacheEntry, DefaultCacheSize),
	}, nil
}

// Get implements RPCServer interface.
func (sm *sessionManager) Context(ctx context.Context, req *empty.Empty) (*mnemosynerpc.ContextResponse, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing metadata in context, access token cannot be retrieved")
	}

	if len(md[mnemosyne.AccessTokenMetadataKey]) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing access token in metadata")
	}

	at := md[mnemosyne.AccessTokenMetadataKey][0]

	res, err := sm.Get(ctx, &mnemosynerpc.GetRequest{AccessToken: at})
	if err != nil {
		return nil, err
	}
	return &mnemosynerpc.ContextResponse{
		Session: res.Session,
	}, nil
}

// Get implements RPCServer interface.
func (sm *sessionManager) Get(ctx context.Context, req *mnemosynerpc.GetRequest) (*mnemosynerpc.GetResponse, error) {
	hs := jump.Sum64(req.AccessToken)
	if node, ok := sm.getNode(hs); ok {
		sklog.Debug(sm.logger, "get request forwarded", "from", sm.addr, "to", node.Addr, "access_token", req.AccessToken)
		return node.Client.Get(ctx, req)
	}
	h := sm.alloc.get(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	var (
		ses *mnemosynerpc.Session
		err error
	)

	entry, ok := sm.readCache(hs)
	if !ok || (!entry.refresh && time.Since(entry.exp) > sm.ttl) {
		if ok {
			if sm.monitor.enabled {
				sm.monitor.cache.refresh.Add(1)
			}
			sm.dataLock.Lock()
			sm.data[hs].refresh = true
			sm.dataLock.Unlock()
		}
		ses, err = h.get(ctx, req)
		if err != nil {
			if grpc.Code(err) == codes.NotFound && ok {
				sm.delCache(hs)
			}
			return nil, err
		}
		sm.putCache(hs, *ses)
	} else {
		ses = &entry.ses
	}

	sklog.Debug(h.logger, "session has been retrieved")

	return &mnemosynerpc.GetResponse{
		Session: ses,
	}, nil
}

// List implements RPCServer interface.
func (sm *sessionManager) List(ctx context.Context, req *mnemosynerpc.ListRequest) (*mnemosynerpc.ListResponse, error) {
	h := sm.alloc.list(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	sessions, err := h.list(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session list has been retrieved")

	return &mnemosynerpc.ListResponse{
		Sessions: sessions,
	}, nil
}

// Start implements RPCServer interface.
func (sm *sessionManager) Start(ctx context.Context, req *mnemosynerpc.StartRequest) (*mnemosynerpc.StartResponse, error) {
	if req.Session == nil {
		return nil, errMissingSession
	}
	if req.Session.AccessToken == "" {
		var err error
		req.Session.AccessToken, err = mnemosyne.RandomAccessToken()
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "access token generation failure: %s", err.Error())
		}
	}
	hs := jump.Sum64(req.Session.AccessToken)
	if node, ok := sm.getNode(hs); ok {
		sklog.Debug(sm.logger, "start request forwarded", "from", sm.addr, "to", node.Addr, "access_token", req.Session.AccessToken)
		return node.Client.Start(ctx, req)
	}
	h := sm.alloc.start(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	ses, err := h.start(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session has been started")

	return &mnemosynerpc.StartResponse{
		Session: ses,
	}, nil
}

// Exists implements RPCServer interface.
func (sm *sessionManager) Exists(ctx context.Context, req *mnemosynerpc.ExistsRequest) (*mnemosynerpc.ExistsResponse, error) {
	h := sm.alloc.exists(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	exists, err := h.exists(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session presence has been checked")

	return &mnemosynerpc.ExistsResponse{
		Exists: exists,
	}, nil
}

// Abandon implements RPCServer interface.
func (sm *sessionManager) Abandon(ctx context.Context, req *mnemosynerpc.AbandonRequest) (*mnemosynerpc.AbandonResponse, error) {
	h := sm.alloc.abandon(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	abandoned, err := h.abandon(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session has been abandoned")

	return &mnemosynerpc.AbandonResponse{
		Abandoned: abandoned,
	}, nil
}

// SetValue implements RPCServer interface.
func (sm *sessionManager) SetValue(ctx context.Context, req *mnemosynerpc.SetValueRequest) (*mnemosynerpc.SetValueResponse, error) {
	h := sm.alloc.setValue(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	bag, err := h.setValue(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session bag value has been set")

	return &mnemosynerpc.SetValueResponse{
		Bag: bag,
	}, nil
}

// Delete implements RPCServer interface.
func (sm *sessionManager) Delete(ctx context.Context, req *mnemosynerpc.DeleteRequest) (*mnemosynerpc.DeleteResponse, error) {
	h := sm.alloc.delete(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	affected, err := h.delete(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session value has been deleted")

	return &mnemosynerpc.DeleteResponse{
		Count: affected,
	}, nil
}

func (sm *sessionManager) cleanup(done chan struct{}) {
	logger := log.NewContext(sm.logger).WithPrefix("module", "cleanup")
	sklog.Info(sm.logger, "cleanup routing started")
InfLoop:
	for {
		select {
		case <-time.After(sm.ttc):
			t := time.Now()
			sklog.Debug(logger, "session cleanup start", "start_at", t.Format(time.RFC3339))
			affected, err := sm.storage.Delete(context.Background(), "", nil, &t)
			if err != nil {
				if sm.monitor.enabled {
					sm.monitor.cleanup.errors.Inc()
				}
				sklog.Error(logger, fmt.Errorf("session cleanup failure: %s", err.Error()), "expire_at_to", t)
				return
			}

			sklog.Debug(logger, "session cleanup success", "count", affected, "elapsed", time.Now().Sub(t))
		case <-done:
			sklog.Info(logger, "cleanup routing terminated")
			break InfLoop
		}
	}
}

func (sm *sessionManager) loggerBackground(ctx context.Context, keyval ...interface{}) log.Logger {
	l := log.NewContext(sm.logger).With(keyval...)
	if md, ok := metadata.FromContext(ctx); ok {
		if rid, ok := md["request_id"]; ok && len(rid) >= 1 {
			l = l.With("request_id", rid[0])
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		l = l.With("peer_address", p.Addr.String())
	}

	return l
}

func (sm *sessionManager) getNode(k uint64) (*cluster.Node, bool) {
	if sm.cluster == nil {
		return nil, false
	}
	if sm.cluster.Len() == 1 {
		return nil, false
	}
	if node, ok := sm.cluster.Get(jump.Hash(k, sm.cluster.Len())); ok {
		if node.Addr != sm.addr {
			if node.Client != nil {
				return node, true
			}
		}
	}
	return nil, false
}

func (sm *sessionManager) putCache(k uint64, ses mnemosynerpc.Session) {
	sm.dataLock.Lock()
	sm.data[k] = &cacheEntry{ses: ses, exp: time.Now().Add(sm.ttl), refresh: false}
	sm.dataLock.Unlock()
}

func (sm *sessionManager) delCache(k uint64) {
	sm.dataLock.Lock()
	delete(sm.data, k)
	sm.dataLock.Unlock()
}

func (sm *sessionManager) readCache(k uint64) (*cacheEntry, bool) {
	sm.dataLock.RLock()
	entry, ok := sm.data[k]
	sm.dataLock.RUnlock()
	if sm.monitor.enabled {
		if ok {
			sm.monitor.cache.hits.Add(1)
		} else {
			sm.monitor.cache.misses.Add(1)
		}
	}
	return entry, ok
}
