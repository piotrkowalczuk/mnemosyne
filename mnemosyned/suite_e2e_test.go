package mnemosyned

import (
	"net"
	"testing"
	"time"

	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/sklog"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
)

type e2eSuite struct {
	listener   net.Listener
	daemon     *Daemon
	client     mnemosynerpc.SessionManagerClient
	clientConn *grpc.ClientConn
}

func (es *e2eSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("e2e suite ignored in short mode")
	}
	var err error
	//logger := sklog.NewHumaneLogger(os.Stdout, sklog.DefaultHTTPFormatter)
	logger := sklog.NewTestLogger(t)

	es.listener = listenTCP(t)
	es.daemon, err = NewDaemon(&DaemonOpts{
		IsTest:          true,
		RPCOptions:      []grpc.ServerOption{},
		RPCListener:     es.listener,
		Storage:         StorageEnginePostgres,
		Logger:          logger,
		PostgresAddress: testPostgresAddress,
		Monitoring:      true,
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
	if err := es.clientConn.Close(); err != nil {
		t.Fatalf("e2e suite client connection unexpected error on close: %s", err.Error())
	}
	if err := es.daemon.Close(); err != nil {
		t.Fatalf("e2e suite daemon unexpected error on close: %s", err.Error())
	}
	if err := es.daemon.storage.TearDown(); err != nil {
		t.Fatalf("e2e suite storage unexpected error on teardown: %s", err.Error())
	} else {
		t.Log("e2e suite storage teardown")
	}
}

func WithE2ESuite(t *testing.T, f func(*e2eSuite)) func() {
	return func() {
		s := &e2eSuite{}
		s.setup(t)

		Reset(func() {
			s.teardown(t)
		})
		f(s)

	}
}
