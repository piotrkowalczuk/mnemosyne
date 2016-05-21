package mnemosyned

import (
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
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

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	d := NewDaemon(&DaemonOpts{
		Environment: EnvironmentTest,
		SessionTTL:  ttl,
		SessionTTC:  ttc,
		RPCListener: l,
		// Use this logger to debug issues
		Logger: sklog.NewHumaneLogger(os.Stdout, sklog.DefaultHTTPFormatter),
		//Logger:                 sklog.NewTestLogger(t),
		StoragePostgresAddress: testPostgresAddress,
	})
	if err := d.Run(); err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	conn, err := grpc.Dial(d.Addr().String(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	m := mnemosyne.New(conn, mnemosyne.MnemosyneOpts{})
	ats := make([]mnemosyne.AccessToken, 0, nb)
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
		_, err := m.Get(ctx, at)
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
