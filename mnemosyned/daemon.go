package mnemosyned

import (
	"database/sql"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"testing"

	"time"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	// EnvironmentProduction ...
	EnvironmentProduction = "prod"
	// EnvironmentTest ...
	EnvironmentTest = "test"
)

// DaemonOpts ...
type DaemonOpts struct {
	Environment            string
	Namespace              string
	Subsystem              string
	SessionTTL             time.Duration
	SessionTTC             time.Duration
	MonitoringEngine       string
	TLS                    bool
	TLSCertFile            string
	TLSKeyFile             string
	StorageEngine          string
	StoragePostgresAddress string
	StoragePostgresTable   string
	Logger                 log.Logger
	RPCOptions             []grpc.ServerOption
	RPCListener            net.Listener
	DebugListener          net.Listener
}

// TestDaemonOpts ...
type TestDaemonOpts struct {
	StoragePostgresAddress string
}

// Daemon ...
type Daemon struct {
	done          chan struct{}
	opts          *DaemonOpts
	monitor       *monitoring
	rpcOptions    []grpc.ServerOption
	storage       Storage
	logger        log.Logger
	rpcListener   net.Listener
	debugListener net.Listener
}

// NewDaemon ...
func NewDaemon(opts *DaemonOpts) *Daemon {
	d := &Daemon{
		done:          make(chan struct{}, 0),
		opts:          opts,
		logger:        opts.Logger,
		rpcOptions:    opts.RPCOptions,
		rpcListener:   opts.RPCListener,
		debugListener: opts.DebugListener,
	}

	if d.opts.SessionTTL == 0 {
		d.opts.SessionTTL = DefaultTTL
	}
	if d.opts.SessionTTC == 0 {
		d.opts.SessionTTC = DefaultTTC
	}
	if d.opts.StorageEngine == "" {
		d.opts.StorageEngine = StorageEnginePostgres
	}
	if d.opts.StoragePostgresTable == "" {
		d.opts.StoragePostgresTable = "session"
	}
	return d
}

// TestDaemon returns address of fully started in-memory daemon and closer to close it.
func TestDaemon(t *testing.T, opts *TestDaemonOpts) (net.Addr, io.Closer) {
	l, err := net.Listen("tcp", "127.0.0.1:0") // any available address
	if err != nil {
		t.Fatalf("mnemosyne daemon tcp listener setup error: %s", err.Error())
	}

	logger := sklog.NewTestLogger(t)
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))

	d := NewDaemon(&DaemonOpts{
		Namespace:              "mnemosyne_test",
		MonitoringEngine:       MonitoringEnginePrometheus,
		Logger:                 logger,
		StoragePostgresAddress: opts.StoragePostgresAddress,
		RPCListener:            l,
	})
	if err := d.Run(); err != nil {
		t.Fatalf("mnemosyne daemon start error: %s", err.Error())
	}

	return d.Addr(), d
}

func (d *Daemon) Run() (err error) {
	if err = d.initMonitoring(); err != nil {
		return
	}
	if err = d.initStorage(); err != nil {
		return
	}

	if d.opts.TLS {
		creds, err := credentials.NewServerTLSFromFile(d.opts.TLSCertFile, d.opts.TLSKeyFile)
		if err != nil {
			return err
		}
		d.rpcOptions = append(d.rpcOptions, grpc.Creds(creds))
	}

	grpclog.SetLogger(sklog.NewGRPCLogger(d.logger))
	gRPCServer := grpc.NewServer(d.rpcOptions...)
	mnemosyneServer := newRPCServer(d.logger, d.storage, d.monitor, d.opts.SessionTTC)
	mnemosyne.RegisterRPCServer(gRPCServer, mnemosyneServer)

	go mnemosyneServer.cleanup(d.done)
	go func() {
		sklog.Info(d.logger, "rpc server is running", "address", d.rpcListener.Addr().String(), "subsystem", d.opts.Subsystem, "namespace", d.opts.Namespace)

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
			sklog.Info(d.logger, "debug server is running", "address", d.debugListener.Addr().String(), "subsystem", d.opts.Subsystem, "namespace", d.opts.Namespace)
			// TODO: implement keep alive
			sklog.Error(d.logger, http.Serve(d.debugListener, nil))
		}()
	}

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

func (d *Daemon) initStorage() (err error) {
	var db *sql.DB

	switch d.opts.StorageEngine {
	case StorageEngineInMemory:
		return errors.New("in memory storage is not implemented yet")
	case StorageEnginePostgres:
		db, err = initPostgres(
			d.opts.StoragePostgresAddress,
			d.logger,
		)
		if err != nil {
			return
		}
		if d.storage, err = initStorage(
			d.opts.Environment,
			newPostgresStorage(d.opts.StoragePostgresTable, db, d.monitor, d.opts.SessionTTL),
			d.logger,
		); err != nil {
			return
		}
		return
	case StorageEngineRedis:
		return errors.New("redis storage is not implemented yet")
	default:
		return errors.New("unknown storage engine")
	}
}

func (d *Daemon) initMonitoring() (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return errors.New("getting hostname failed")
	}

	switch d.opts.MonitoringEngine {
	case "":
		d.monitor = &monitoring{}
		return
	case MonitoringEnginePrometheus:
		d.monitor = initPrometheus(d.opts.Namespace, d.opts.Subsystem, prometheus.Labels{"server": hostname})
		return
	default:
		return errors.New("unknown monitoring engine")
	}
}
