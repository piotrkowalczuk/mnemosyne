package mnemosyned

import (
	"net"
	"testing"
	"time"

	"go.uber.org/zap"

	"context"

	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage/storagemock"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"google.golang.org/grpc"
)

type integrationSuite struct {
	listener      net.Listener
	server        *grpc.Server
	service       mnemosynerpc.SessionManagerClient
	serviceConn   *grpc.ClientConn
	serviceServer mnemosynerpc.SessionManagerServer
	store         *storagemock.Storage
}

func (is *integrationSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("integration suite ignored in short mode")
	}

	var err error

	is.store = &storagemock.Storage{}
	is.listener = listenTCP(t)
	is.server = grpc.NewServer()
	is.serviceServer, err = newSessionManager(sessionManagerOpts{
		addr:    is.listener.Addr().String(),
		logger:  zap.L(),
		storage: is.store,
		cluster: is.initCluster(t),
		ttc:     storage.DefaultTTC,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	mnemosynerpc.RegisterSessionManagerServer(is.server, is.serviceServer)

	go is.server.Serve(is.listener)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	is.serviceConn, err = grpc.DialContext(
		ctx,
		is.listener.Addr().String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
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
	csr, err := cluster.New(cluster.Opts{
		Listen: is.listener.Addr().String(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return csr
}
