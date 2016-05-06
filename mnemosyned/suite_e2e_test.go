package mnemosyned

import (
	"net"
	"testing"
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
)

type e2eSuite struct {
	listener   net.Listener
	daemon     *Daemon
	client     mnemosyne.RPCClient
	clientConn *grpc.ClientConn
}

func (es *e2eSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("e2e suite ignored in short mode")
	}

	logger := sklog.NewTestLogger(t)

	es.listener = listenTCP(t)
	es.daemon = NewDaemon(&DaemonOpts{
		Namespace:              "mnemosyne_e2e",
		Subsystem:              "mnemosyne",
		RPCOptions:             []grpc.ServerOption{},
		RPCListener:            es.listener,
		StorageEngine:          StorageEnginePostgres,
		Logger:                 logger,
		StoragePostgresAddress: testPostgresAddress,
	})

	var err error
	if err = es.daemon.Run(); err != nil {
		t.Fatalf("unexpected deamon run error: %s", err.Error())
	} else {
		t.Logf("test daemon started")
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

	es.client = mnemosyne.NewRPCClient(es.clientConn)
}

func (es *e2eSuite) teardown(t *testing.T) {
	if err := es.clientConn.Close(); err != nil {
		t.Fatalf("e2e suite client connection unexpected error on close: %s", err.Error())
	}
	if err := es.daemon.storage.TearDown(); err != nil {
		t.Fatalf("e2e suite storage unexpected error on teardown: %s", err.Error())
	} else {
		t.Logf("e2e suite storage teardown")
	}
	if err := es.daemon.Close(); err != nil {
		t.Fatalf("e2e suite daemon unexpected error on close: %s", err.Error())
	}
}

func WithE2ESuite(t *testing.T, f func(*e2eSuite)) func() {
	return func() {
		s := &e2eSuite{}
		s.setup(t)

		state, err := s.clientConn.State()

		So(err, ShouldBeNil)
		So(state, ShouldNotEqual, grpc.Shutdown)
		So(state, ShouldNotEqual, grpc.TransientFailure)
		Reset(func() {
			s.teardown(t)
		})
		f(s)

	}
}
