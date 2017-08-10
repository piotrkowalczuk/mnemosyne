package mnemosyned

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
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

	rl := listener(t)
	dl := listener(t)

	d, err := NewDaemon(&DaemonOpts{
		IsTest:        true,
		SessionTTL:    ttl,
		SessionTTC:    ttc,
		RPCListener:   rl,
		DebugListener: dl,
		// Use this logger to debug issues
		//Logger: sklog.NewHumaneLogger(os.Stdout, sklog.DefaultHTTPFormatter),
		Logger:          zap.L(),
		PostgresAddress: testPostgresAddress,
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
		t.Fatalf("unexpected mnemosyne instatiation error: %s", err.Error())
	}
	defer m.Close()

	ats := make([]string, 0, nb)
	for i := 0; i < nb; i++ {

		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		ses, err := m.Start(ctx, strconv.Itoa(i), "daemon test client", nil)
		if err != nil {
			t.Errorf("session could not be started: %s", err.Error())
			return
		}
		t.Logf("session created, it expires at: %s", time.Unix(ses.ExpireAt.Seconds, int64(ses.ExpireAt.Nanos)).Format(time.RFC3339))
		ats = append(ats, ses.AccessToken)
	}

	// BUG: this assertion can fail on travis because of cpu lag.
	<-time.After(ttl + ttc + ttc)

	for i, at := range ats {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		_, err := m.Get(ctx, string(at))
		if err == nil {
			t.Errorf("%d: missing error", i)
			return
		}
		if grpc.Code(err) != codes.NotFound {
			t.Errorf("%d: wrong error code, expected %d but got %d", i, codes.NotFound, grpc.Code(err))
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
		t.Error("address should not be empty")
	}
	if err := closer.Close(); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestDaemon_Cluster(t *testing.T) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zap.WarnLevel)
	l, _ := cfg.Build()

	l1 := listener(t)
	l2 := listener(t)
	l3 := listener(t)

	d1, err := NewDaemon(&DaemonOpts{
		IsTest:            true,
		RPCListener:       l1,
		Logger:            l,
		PostgresAddress:   testPostgresAddress,
		ClusterListenAddr: l1.Addr().String(),
		ClusterSeeds: []string{
			l1.Addr().String(),
			l2.Addr().String(),
			l3.Addr().String(),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if err := d1.Run(); err != nil {
		t.Fatalf("mnemosyne daemon 1 start error: %s", err.Error())
	}
	defer d1.Close()

	d2, err := NewDaemon(&DaemonOpts{
		IsTest:            true,
		RPCListener:       l2,
		Logger:            l,
		PostgresAddress:   testPostgresAddress,
		ClusterListenAddr: l2.Addr().String(),
		ClusterSeeds: []string{
			l1.Addr().String(),
			l2.Addr().String(),
			l3.Addr().String(),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if err := d2.Run(); err != nil {
		t.Fatalf("mnemosyne daemon 2 start error: %s", err.Error())
	}
	defer d2.Close()

	d3, err := NewDaemon(&DaemonOpts{
		IsTest:            true,
		RPCListener:       l3,
		Logger:            l,
		PostgresAddress:   testPostgresAddress,
		ClusterListenAddr: l3.Addr().String(),
		ClusterSeeds: []string{
			l1.Addr().String(),
			l2.Addr().String(),
			l3.Addr().String(),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if err := d3.Run(); err != nil {
		t.Fatalf("mnemosyne daemon 3 start error: %s", err.Error())
	}
	defer d3.Close()

	c1, m1 := connect(t, l1)
	defer c1.Close()
	c2, m2 := connect(t, l2)
	defer c2.Close()
	c3, m3 := connect(t, l3)
	defer c3.Close()

	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel1()
	r1, err := m1.Start(ctx1, &mnemosynerpc.StartRequest{
		Session: &mnemosynerpc.Session{
			SubjectId:     "1",
			SubjectClient: "mnemosyne-test",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	r2, err := m2.Start(ctx2, &mnemosynerpc.StartRequest{
		Session: &mnemosynerpc.Session{
			SubjectId:     "2",
			SubjectClient: "mnemosyne-test",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	ctx3, cancel3 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel3()
	r3, err := m3.Start(ctx3, &mnemosynerpc.StartRequest{
		Session: &mnemosynerpc.Session{
			SubjectId:     "3",
			SubjectClient: "mnemosyne-test",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	clients := map[string]mnemosynerpc.SessionManagerClient{
		r1.Session.AccessToken: m1,
		r2.Session.AccessToken: m2,
		r3.Session.AccessToken: m3,
	}

	var i, j int
	for at, c := range clients {
		i++

		t.Run(fmt.Sprintf("client_%d", i), func(t *testing.T) {
			for range clients {
				j++

				t.Run(fmt.Sprintf("server_%d/get", j), func(t *testing.T) {
					ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
					defer cancel()
					res, err := c.Get(ctx, &mnemosynerpc.GetRequest{
						AccessToken: at,
					})
					if err != nil {
						t.Fatal(err)
					}
					if res.Session.AccessToken != at {
						t.Fatal("wrong access token")
					}
				})
				t.Run(fmt.Sprintf("server_%d/exists", j), func(t *testing.T) {
					ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
					defer cancel()

					res, err := c.Exists(ctx, &mnemosynerpc.ExistsRequest{
						AccessToken: at,
					})
					if err != nil {
						t.Fatal(err)
					}
					if !res.Exists {
						t.Fatal("should exists")
					}
				})
			}
		})
	}
}

func listener(t testing.TB) net.Listener {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func connect(t testing.TB, l net.Listener) (*grpc.ClientConn, mnemosynerpc.SessionManagerClient) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := grpc.DialContext(ctx, l.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial error: %s", err.Error())
	}

	return c, mnemosynerpc.NewSessionManagerClient(c)
}
