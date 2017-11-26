package mnemosyned

import (
	"fmt"
	"net"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
)

type e2eSuites []*e2eSuite

func (es *e2eSuites) setup(t *testing.T, factor int) {
	if testing.Short() {
		t.Skip("e2e suite ignored in short mode")
	}
	t.Logf("e2e suite setup start for factor %d", factor)

	listeners := make([]net.Listener, 0, factor)
	for i := 0; i < factor; i++ {
		listeners = append(listeners, listenTCP(t))
	}
	for i, l := range listeners {
		var seeds []string
		for _, s := range listeners {
			if s.Addr().String() != l.Addr().String() {
				seeds = append(seeds, s.Addr().String())
			}
		}
		s := &e2eSuite{
			listener: l,
			seeds:    seeds,
		}
		s.setup(t, i)
		*es = append(*es, s)
	}
	t.Logf("e2e suite setup finish for factor %d", factor)
}

func (es e2eSuites) teardown(t *testing.T) {
	for _, s := range es {
		s.teardown(t)
	}
}

type e2eSuite struct {
	listener   net.Listener
	seeds      []string
	daemon     *Daemon
	client     mnemosynerpc.SessionManagerClient
	clientConn *grpc.ClientConn
}

func (es *e2eSuite) setup(t *testing.T, i int) {
	var err error

	if es.listener == nil {
		es.listener = listenTCP(t)
	}
	es.daemon, err = NewDaemon(&DaemonOpts{
		IsTest:            true,
		RPCOptions:        []grpc.ServerOption{},
		RPCListener:       es.listener,
		Storage:           storage.EnginePostgres,
		Logger:            zap.L(),
		PostgresAddress:   testPostgresAddress,
		PostgresSchema:    fmt.Sprintf("mnemosyne_test_%d", i),
		Monitoring:        true,
		ClusterListenAddr: es.listener.Addr().String(),
		ClusterSeeds:      es.seeds,
	})
	if err != nil {
		t.Fatalf("unexpected deamon instantiation error: %s", err.Error())
	}
	if err = es.daemon.Run(); err != nil {
		t.Fatalf("unexpected deamon run error: %s", err.Error())
	} else {
		t.Log("test daemon started")
	}

	es.clientConn, err = grpc.Dial(
		es.listener.Addr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
	)
	if err != nil {
		t.Fatalf("unexpected client conn error: %s", err.Error())
	}

	es.client = mnemosynerpc.NewSessionManagerClient(es.clientConn)
}

func (es *e2eSuite) teardown(t *testing.T) {
	t.Helper()

	if err := es.clientConn.Close(); err != nil {
		t.Fatalf("e2e suite client connection unexpected error on close: %s", err.Error())
	}
	if err := es.daemon.Close(); err != nil {
		t.Fatalf("e2e suite daemon unexpected error on close: %s", err.Error())
	}
}

func WithE2ESuite(t *testing.T, f func(*e2eSuite)) func() {
	return func() {
		s := e2eSuites{}
		s.setup(t, 1)

		Reset(func() {
			s.teardown(t)
		})

		f(s[0])
	}
}

func WithE2ESuites(t *testing.T, factor int, f func(e2eSuites)) func() {
	return func() {
		s := e2eSuites{}
		s.setup(t, factor)

		Reset(func() {
			s.teardown(t)
		})

		f(s)
	}
}
