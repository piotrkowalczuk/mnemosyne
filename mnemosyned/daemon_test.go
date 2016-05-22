package mnemosyned

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestDaemon_Run(t *testing.T) {
	if testing.Short() {
		t.Skip("this test takes to long to run it in short mode")
	}

	ttl := 5 * time.Second
	ttc := 1 * time.Second
	nb := 10

	rl, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	dl, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	d, err := NewDaemon(&DaemonOpts{
		Environment:   EnvironmentTest,
		SessionTTL:    ttl,
		SessionTTC:    ttc,
		RPCListener:   rl,
		DebugListener: dl,
		// Use this logger to debug issues
		//Logger: sklog.NewHumaneLogger(os.Stdout, sklog.DefaultHTTPFormatter),
		Logger:                 sklog.NewTestLogger(t),
		StoragePostgresAddress: testPostgresAddress,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = d.Run(); err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	m, err := mnemosyne.New(mnemosyne.MnemosyneOpts{
		Addresses: []string{d.Addr().String()},
	})
	if err != nil {
		t.Fatal("unexpected mnemosyne instatiation error: %s", err.Error())
	}
	defer m.Close()

	ats := make([]mnemosynerpc.AccessToken, 0, nb)
	for i := 0; i < nb; i++ {

		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		ses, err := m.Start(ctx, strconv.Itoa(i), "daemon test client", nil)
		if err != nil {
			t.Errorf("session could not be started:", err)
			return
		}
		t.Logf("session created, it expires at: %s", time.Unix(ses.ExpireAt.Seconds, int64(ses.ExpireAt.Nanos)).Format(time.RFC3339))
		ats = append(ats, *ses.AccessToken)
	}

	// BUG: this assertion can fail on travis because of cpu lag.
	<-time.After(ttl + ttc + ttc)

	for i, at := range ats {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		_, err := m.Get(ctx, at.Encode())
		if err == nil {
			t.Error("%d: missing error", i)
			return
		}
		if grpc.Code(err) != codes.NotFound {
			t.Errorf("%d: wrong error code, expected %s", i, codes.NotFound, grpc.Code(err))
			return
		}

		t.Logf("%d: as expected, session does not exists anymore", i)
	}
}

func TestTestDaemon(t *testing.T) {
	addr, closer := TestDaemon(t, TestDaemonOpts{
		StoragePostgresAddress: testPostgresAddress,
	})
	if addr.String() == "" {
		t.Errorf("address should not be empty")
	}
	if err := closer.Close(); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}
