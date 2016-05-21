package mnemosyned

import (
	"net"
	"strconv"
	"testing"
	"time"

	"sync"

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
		//Logger:                 sklog.NewHumaneLogger(os.Stdout, sklog.DefaultHTTPFormatter),
		Logger:                 sklog.NewTestLogger(t),
		StoragePostgresAddress: testPostgresAddress,
	})
	if err := d.Run(); err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	conn, err := grpc.Dial(d.Addr().String(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		t.Fatal(err)
	}

	m := mnemosyne.New(conn, mnemosyne.MnemosyneOpts{})
	ats := make([]mnemosyne.AccessToken, 0, nb)
	wg := &sync.WaitGroup{}
	for i := 0; i < nb; i++ {
		go func(i int, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()

			ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
			ses, err := m.Start(ctx, strconv.Itoa(i), "daemon test client", nil)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("session created, it expires at: %s", time.Unix(ses.ExpireAt.Seconds, int64(ses.ExpireAt.Nanos)).Format(time.RFC3339))
			ats = append(ats, *ses.AccessToken)
		}(i, wg)
	}
	wg.Wait()

	<-time.After(ttl + ttc + ttc)

	for _, at := range ats {
		go func(at mnemosyne.AccessToken, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()

			ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
			_, err := m.Get(ctx, at)
			if err == nil {
				t.Error("missing error")
				return
			}
			if grpc.Code(err) != codes.NotFound {
				t.Errorf("wrong error code, expected %s", codes.NotFound, grpc.Code(err))
			}
		}(at, wg)
	}
	wg.Wait()
}
