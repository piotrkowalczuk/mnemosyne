package mnemosyned

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"strings"
	"testing"
	"time"

	"context"

	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/service/postgres"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	storagepq "github.com/piotrkowalczuk/mnemosyne/internal/storage/postgres"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/promgrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const subsystem = "mnemosyned"

// DaemonOpts it is constructor argument that can be passed to
// the NewDaemon constructor function.
type DaemonOpts struct {
	Version           string
	IsTest            bool
	SessionTTL        time.Duration
	SessionTTC        time.Duration
	TLS               bool
	TLSCertFile       string
	TLSKeyFile        string
	Storage           string
	PostgresAddress   string
	PostgresTable     string
	PostgresSchema    string
	Logger            *zap.Logger
	RPCOptions        []grpc.ServerOption
	RPCListener       net.Listener
	DebugListener     net.Listener
	ClusterListenAddr string
	ClusterSeeds      []string
}

// TestDaemonOpts set of options that are used with TestDaemon instance.
type TestDaemonOpts struct {
	StoragePostgresAddress string
}

// Daemon represents single daemon instance that can be run.
type Daemon struct {
	opts          *DaemonOpts
	done          chan struct{}
	serverOptions []grpc.ServerOption
	clientOptions []grpc.DialOption
	postgres      *sql.DB
	logger        *zap.Logger
	server        *grpc.Server
	storage       storage.Storage
	rpcListener   net.Listener
	debugListener net.Listener
}

// NewDaemon allocates new daemon instance using given options.
func NewDaemon(opts *DaemonOpts) (*Daemon, error) {
	d := &Daemon{
		done:          make(chan struct{}),
		opts:          opts,
		logger:        opts.Logger,
		serverOptions: opts.RPCOptions,
		rpcListener:   opts.RPCListener,
		debugListener: opts.DebugListener,
	}

	if err := d.setPostgresConnectionParameters(); err != nil {
		return nil, err
	}
	if d.opts.SessionTTL == 0 {
		d.opts.SessionTTL = storage.DefaultTTL
	}
	if d.opts.SessionTTC == 0 {
		d.opts.SessionTTC = storage.DefaultTTC
	}
	if d.opts.Storage == "" {
		d.opts.Storage = storage.EnginePostgres
	}
	if d.opts.PostgresTable == "" {
		d.opts.PostgresTable = "session"
	}
	if d.opts.PostgresSchema == "" {
		d.opts.PostgresSchema = "mnemosyne"
	}

	return d, nil
}

// TestDaemon returns address of fully started in-memory daemon and closer to close it.
func TestDaemon(t *testing.T, opts TestDaemonOpts) (net.Addr, io.Closer) {
	l, err := net.Listen("tcp", "127.0.0.1:0") // any available address
	if err != nil {
		t.Fatalf("mnemosyne daemon tcp listener setup error: %s", err.Error())
	}

	d, err := NewDaemon(&DaemonOpts{
		IsTest:            true,
		ClusterListenAddr: l.Addr().String(),
		Logger:            zap.L(),
		PostgresAddress:   opts.StoragePostgresAddress,
		PostgresTable:     "session",
		PostgresSchema:    "mnemosyne",
		RPCListener:       l,
	})
	if err != nil {
		t.Fatalf("mnemosyne daemon cannot be instantiated: %s", err.Error())
	}
	if err := d.Run(); err != nil {
		t.Fatalf("mnemosyne daemon start error: %s", err.Error())
	}

	return d.Addr(), d
}

// Run starts daemon and all services within.
func (d *Daemon) Run() (err error) {
	var (
		cl *cluster.Cluster
	)
	if cl, err = initCluster(d.logger, d.opts.ClusterListenAddr, d.opts.ClusterSeeds...); err != nil {
		return
	}
	if err = d.initStorage(d.logger, d.opts.PostgresTable, d.opts.PostgresSchema); err != nil {
		return
	}

	interceptor := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})

	d.clientOptions = []grpc.DialOption{
		grpc.WithUserAgent(fmt.Sprintf("%s:%s", subsystem, d.opts.Version)),
		grpc.WithStatsHandler(interceptor),
		grpc.WithDialer(interceptor.Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("tcp", addr, timeout)
		})),
		grpc.WithUnaryInterceptor(interceptor.UnaryClient()),
		grpc.WithStreamInterceptor(interceptor.StreamClient()),
	}
	d.serverOptions = []grpc.ServerOption{
		grpc.StatsHandler(interceptor),
		grpc.UnaryInterceptor(unaryServerInterceptors(
			errorInterceptor(d.logger),
			interceptor.UnaryServer(),
		)),
	}
	if d.opts.TLS {
		servCreds, err := credentials.NewServerTLSFromFile(d.opts.TLSCertFile, d.opts.TLSKeyFile)
		if err != nil {
			return err
		}
		d.serverOptions = append(d.serverOptions, grpc.Creds(servCreds))

		clientCreds, err := credentials.NewClientTLSFromFile(d.opts.TLSCertFile, "")
		if err != nil {
			return err
		}
		d.clientOptions = append(d.clientOptions, grpc.WithTransportCredentials(clientCreds))
	} else {
		d.clientOptions = append(d.clientOptions, grpc.WithInsecure())
	}

	d.server = grpc.NewServer(d.serverOptions...)

	cache := cache.New(5*time.Second, subsystem)
	mnemosyneServer, err := newSessionManager(sessionManagerOpts{
		addr:    d.opts.ClusterListenAddr,
		cluster: cl,
		logger:  d.logger,
		storage: d.storage,
		ttc:     d.opts.SessionTTC,
		cache:   cache,
	})
	if err != nil {
		return err
	}

	mnemosynerpc.RegisterSessionManagerServer(d.server, mnemosyneServer)
	grpc_health_v1.RegisterHealthServer(d.server, health.NewServer())

	if !d.opts.IsTest {
		prometheus.DefaultRegisterer.Register(d.storage.(storage.InstrumentedStorage))
		prometheus.DefaultRegisterer.Register(cache)
		prometheus.DefaultRegisterer.Register(mnemosyneServer)
		prometheus.DefaultRegisterer.Register(interceptor)
		promgrpc.RegisterInterceptor(d.server, interceptor)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = cl.Connect(ctx, d.clientOptions...); err != nil {
		return err
	}

	go func() {
		d.logger.Info("rpc server is running", zap.String("address", d.rpcListener.Addr().String()))

		if err := d.server.Serve(d.rpcListener); err != nil {
			if err == grpc.ErrServerStopped {
				d.logger.Info("grpc server has been stopped")
				return
			}

			if !strings.Contains(err.Error(), "use of closed network connection") {
				d.logger.Error("rpc server failure", zap.Error(err))
			}
		}
	}()

	if d.debugListener != nil {
		go func() {
			d.logger.Info("debug server is running", zap.String("address", d.debugListener.Addr().String()))

			mux := http.NewServeMux()
			mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
			mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
			mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
			mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
			mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
			mux.Handle("/healthz", &livenessHandler{
				livenessResponse: livenessResponse{
					Version: d.opts.Version,
				},
				logger: d.logger,
			})
			mux.Handle("/healthr", &readinessHandler{
				livenessResponse: livenessResponse{
					Version: d.opts.Version,
				},
				logger:   d.logger,
				postgres: d.postgres,
				cluster:  cl,
			})
			if err := http.Serve(d.debugListener, mux); err != nil {
				d.logger.Error("debug server failure", zap.Error(err))
			}
		}()
	}

	go mnemosyneServer.cleanup(d.done)

	return
}

// Close implements io.Closer interface.
func (d *Daemon) Close() (err error) {
	d.done <- struct{}{}
	d.server.GracefulStop()
	if d.postgres != nil {
		if err = d.postgres.Close(); err != nil {
			return
		}
	}
	if d.debugListener != nil {
		err = d.debugListener.Close()
	}

	return
}

// Addr returns net.Addr that rpc service is listening on.
func (d *Daemon) Addr() net.Addr {
	return d.rpcListener.Addr()
}

func (d *Daemon) initStorage(l *zap.Logger, table, schema string) (err error) {
	switch d.opts.Storage {
	case storage.EngineInMemory:
		return errors.New("in memory storage is not implemented yet")
	case storage.EnginePostgres:
		d.postgres, err = postgres.Init(
			d.opts.PostgresAddress,
			postgres.Opts{
				Logger: d.logger,
			},
		)
		if err != nil {
			return
		}
		if d.storage, err = storage.Init(storagepq.NewStorage(storagepq.StorageOpts{
			Namespace: subsystem,
			Schema:    schema,
			Table:     table,
			Conn:      d.postgres,
			TTL:       d.opts.SessionTTL,
		}), d.opts.IsTest); err != nil {
			return
		}

		l.Info("postgres storage initialized", zap.String("schema", schema), zap.String("table", table))
		return
	case storage.EngineRedis:
		return errors.New("redis storage is not implemented yet")
	default:
		return errors.New("unknown storage engine")
	}
}

func (d *Daemon) setPostgresConnectionParameters() error {
	u, err := url.Parse(d.opts.PostgresAddress)
	if err != nil {
		return err
	}
	v := u.Query()
	v.Set("timezone", "utc")
	u.RawQuery = v.Encode()
	d.opts.PostgresAddress = u.String()
	return nil
}
