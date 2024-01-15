package logger_test

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"

	"github.com/piotrkowalczuk/mnemosyne/internal/service/logger"
)

func TestInit(t *testing.T) {
	testInit(t, logger.Opts{Environment: "production"})
	testInit(t, logger.Opts{Environment: "development"})
	testInit(t, logger.Opts{Environment: "stackdriver"})
	testInit(t, logger.Opts{Level: "info"})
}

func testInit(t *testing.T, opts logger.Opts) {
	defer func() {
		t.Log("recover", recover())
	}()

	t.Helper()

	l, err := logger.Init(opts)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	ctx := logger.Ctx(peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.IPv4(10, 0, 0, 1),
			Port: 1000,
		},
	}), &grpc.UnaryServerInfo{FullMethod: "fullMethod"}, codes.OK)
	l.Info("info message", ctx)
	l.Debug("debug message", ctx)
	l.Error("error message", ctx)
	l.Warn("warning message", ctx)

	l.DPanic("dpanic message", ctx)
	l.Panic("panic message", ctx)
}
