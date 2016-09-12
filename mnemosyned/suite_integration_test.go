package mnemosyned

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/sklog"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

type integrationSuite struct {
	logger        log.Logger
	listener      net.Listener
	server        *grpc.Server
	service       mnemosynerpc.SessionManagerClient
	serviceConn   *grpc.ClientConn
	serviceServer mnemosynerpc.SessionManagerServer
	store         *mockStorage
}

func (is *integrationSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("integration suite ignored in short mode")
	}

	var err error

	logger := sklog.NewTestLogger(t)
	//logger := sklog.NewHumaneLogger(os.Stdout, sklog.DefaultHTTPFormatter)
	monitor := initPrometheus("mnemosyne_test", false, stdprometheus.Labels{"server": "test"})

	is.store = &mockStorage{}
	is.listener = listenTCP(t)
	is.server = grpc.NewServer()
	is.serviceServer, err = newSessionManager(sessionManagerOpts{
		addr:       is.listener.Addr().String(),
		logger:     logger,
		storage:    is.store,
		monitoring: monitor,
		cluster:    is.initCluster(t),
		ttc:        DefaultTTC,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	mnemosynerpc.RegisterSessionManagerServer(is.server, is.serviceServer)

	go is.server.Serve(is.listener)
	is.serviceConn, err = grpc.Dial(
		is.listener.Addr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	is.service = mnemosynerpc.NewSessionManagerClient(is.serviceConn)
}

func (is *integrationSuite) teardown(t *testing.T) {
	if err := is.serviceConn.Close(); err != nil {
		t.Errorf("integration suite unexpected error on service close: %s", err.Error())
	}
	if err := is.listener.Close(); err != nil {
		t.Errorf("integration suite unexpected error on listener close: %s", err.Error())
	}
}

func (is *integrationSuite) initCluster(t *testing.T) *cluster.Cluster {
	csr, err := cluster.New(is.listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	return csr
}
