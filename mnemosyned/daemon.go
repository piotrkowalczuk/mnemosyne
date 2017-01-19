package mnemosyned

import (
	"database/sql"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/promgrpc"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// DaemonOpts it is constructor argument that can be passed to
// the NewDaemon constructor function.
type DaemonOpts struct {
	IsTest            bool
	SessionTTL        time.Duration
	SessionTTC        time.Duration
	Monitoring        bool
	TLS               bool
	TLSCertFile       string
	TLSKeyFile        string
	Storage           string
	PostgresAddress   string
	PostgresTable     string
	PostgresSchema    string
	Logger            log.Logger
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
	done          chan struct{}
	opts          *DaemonOpts
	monitor       *monitoring
	serverOptions []grpc.ServerOption
	clientOptions []grpc.DialOption
	postgres      *sql.DB
	storage       storage
	logger        log.Logger
	rpcListener   net.Listener
	debugListener net.Listener
}

// NewDaemon allocates new daemon instance using given options.
func NewDaemon(opts *DaemonOpts) (*Daemon, error) {
	d := &Daemon{
		done:          make(chan struct{}, 0),
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
		d.opts.SessionTTL = DefaultTTL
	}
	if d.opts.SessionTTC == 0 {
		d.opts.SessionTTC = DefaultTTC
	}
	if d.opts.Storage == "" {
		d.opts.Storage = StorageEnginePostgres
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

	logger := sklog.NewTestLogger(t)
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))

	d, err := NewDaemon(&DaemonOpts{
		IsTest:            true,
		Monitoring:        false,
		ClusterListenAddr: l.Addr().String(),
		Logger:            logger,
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
		cluster *cluster.Cluster
	)
	if cluster, err = initCluster(d.logger, d.opts.ClusterListenAddr, d.opts.ClusterSeeds...); err != nil {
		return
	}
	if err = d.initMonitoring(); err != nil {
		return
	}
	if err = d.initStorage(d.logger, d.opts.PostgresTable, d.opts.PostgresSchema); err != nil {
		return
	}

	if d.opts.TLS {
		creds, err := credentials.NewServerTLSFromFile(d.opts.TLSCertFile, d.opts.TLSKeyFile)
		if err != nil {
			return err
		}
		d.serverOptions = append(d.serverOptions, grpc.Creds(creds))
		d.clientOptions = append(d.clientOptions, grpc.WithTransportCredentials(creds))
	} else {
		d.clientOptions = append(d.clientOptions, grpc.WithInsecure())
	}

	interceptor := promgrpc.NewInterceptor()
	grpclog.SetLogger(sklog.NewGRPCLogger(d.logger))
	gRPCServer := grpc.NewServer(append(
		d.serverOptions,
		// No stream endpoint available at the moment.
		grpc.UnaryInterceptor(unaryServerInterceptors(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				res, err := handler(ctx, req)

				if err != nil && grpc.Code(err) != codes.OK {
					switch err {
					case errMissingSession, errMissingAccessToken, errMissingSubjectID, errSessionNotFound:
						return nil, err
					}

					code := grpc.Code(err)
					switch code {
					case codes.Unknown:
						sklog.Debug(d.logger, "request handled successfully", "handler", info.FullMethod)
						return nil, grpc.Errorf(codes.Internal, "mnemosyned: %s", grpc.ErrorDesc(err))
					default:
						sklog.Error(d.logger, err, "handler", info.FullMethod)
						return nil, grpc.Errorf(code, "mnemosyned: %s", grpc.ErrorDesc(err))
					}
				}
				return res, err
			},
			interceptor.UnaryServer(),
		)),
	)...)

	mnemosyneServer, err := newSessionManager(sessionManagerOpts{
		addr:       d.opts.ClusterListenAddr,
		cluster:    cluster,
		logger:     d.logger,
		storage:    d.storage,
		monitoring: d.monitor,
		ttc:        d.opts.SessionTTC,
	})
	if err != nil {
		return err
	}
	mnemosynerpc.RegisterSessionManagerServer(gRPCServer, mnemosyneServer)
	promgrpc.RegisterInterceptor(gRPCServer, interceptor)

	if err = cluster.Connect(d.clientOptions...); err != nil {
		return err
	}

	go func() {
		sklog.Info(d.logger, "rpc server is running", "address", d.rpcListener.Addr().String())

		if err := gRPCServer.Serve(d.rpcListener); err != nil {
			if err == grpc.ErrServerStopped {
				sklog.Info(d.logger, "grpc server has been stoped")
				return
			}

			sklog.Error(d.logger, err)
		}
	}()

	if d.debugListener != nil {
		go func() {
			sklog.Info(d.logger, "debug server is running", "address", d.debugListener.Addr().String())
			// TODO: implement keep alive

			mux := http.NewServeMux()
			mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
			mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
			mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
			mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
			if d.opts.Monitoring {
				mux.Handle("/metrics", prometheus.Handler())
			}
			mux.Handle("/health", &healthHandler{
				postgres: d.postgres,
			})
			sklog.Error(d.logger, http.Serve(d.debugListener, mux))
		}()
	}

	go mnemosyneServer.cleanup(d.done)

	return
}

// Close implements io.Closer interface.
func (d *Daemon) Close() (err error) {
	d.done <- struct{}{}
	if err = d.rpcListener.Close(); err != nil {
		return
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

func (d *Daemon) initStorage(l log.Logger, table, schema string) (err error) {
	switch d.opts.Storage {
	case StorageEngineInMemory:
		return errors.New("in memory storage is not implemented yet")
	case StorageEnginePostgres:
		d.postgres, err = initPostgres(
			d.opts.PostgresAddress,
			d.logger,
		)
		if err != nil {
			return
		}
		if d.storage, err = initStorage(
			d.opts.IsTest,
			newPostgresStorage(table, schema, d.postgres, d.monitor, d.opts.SessionTTL),
			d.logger,
		); err != nil {
			return
		}

		sklog.Info(l, "postgres storage initialized", "schema", schema, "table", table)
		return
	case StorageEngineRedis:
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

func (d *Daemon) initMonitoring() (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return errors.New("getting hostname failed")
	}

	d.monitor = initPrometheus("mnemosyne", d.opts.Monitoring, prometheus.Labels{"server": hostname})
	return
}
