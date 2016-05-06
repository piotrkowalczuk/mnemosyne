package mnemosyned

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

type integrationSuite struct {
	logger        log.Logger
	listener      net.Listener
	server        *grpc.Server
	service       mnemosyne.RPCClient
	serviceConn   *grpc.ClientConn
	serviceServer mnemosyne.RPCServer
	store         *storageMock
}

func (is *integrationSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("integration suite ignored in short mode")
	}

	var err error

	logger := sklog.NewTestLogger(t)
	monitor := initPrometheus("mnemosyne_test", "mnemosyne", stdprometheus.Labels{"server": "test"})

	is.store = &storageMock{}
	is.listener = listenTCP(t)
	is.server = grpc.NewServer()
	is.serviceServer = newRPCServer(logger, is.store, monitor)

	mnemosyne.RegisterRPCServer(is.server, is.serviceServer)

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
	is.service = mnemosyne.NewRPCClient(is.serviceConn)
}

func (is *integrationSuite) teardown(t *testing.T) {
	if err := is.serviceConn.Close(); err != nil {
		t.Errorf("integration suite unexpected error on service close: %s", err.Error())
	}
	if err := is.listener.Close(); err != nil {
		t.Errorf("integration suite unexpected error on listener close: %s", err.Error())
	}
}
