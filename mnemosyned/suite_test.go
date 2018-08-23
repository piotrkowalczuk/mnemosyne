package mnemosyned

import (
	"flag"
	"fmt"
	"net"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/status"
)

var (
	testPostgresAddress string
)

func TestMain(m *testing.M) {
	flag.StringVar(&testPostgresAddress, "postgres.address", getStringEnvOr("MNEMOSYNED_POSTGRES_ADDRESS", "postgres://localhost/test?sslmode=disable"), "")
	flag.Parse()

	os.Exit(m.Run())
}

func getStringEnvOr(env, or string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return or
}

func listenTCP(t *testing.T) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0") // any available address
	if err != nil {
		t.Fatalf("net.Listen tcp :0: %s", err.Error())
	}
	return l
}

func ShouldBeValidStartResponse(actual interface{}, expected ...interface{}) (s string) {
	if len(expected) != 1 {
		return fmt.Sprintf("This assertion requires exactly 1 comparison values (you provided %d).", len(expected))
	}

	sr, ok := actual.(*mnemosynerpc.StartResponse)
	if !ok {
		return "The given value must be *StartResponse."
	}
	if s = ShouldBeValidSession(sr.Session); s != "" {
		return
	}
	return
}

func ShouldBeValidGetResponse(actual interface{}, expected ...interface{}) (s string) {
	if len(expected) != 1 {
		return fmt.Sprintf("This assertion requires exactly 1 comparison values (you provided %d).", len(expected))
	}

	gr, ok := actual.(*mnemosynerpc.GetResponse)
	if !ok {
		return "The given value must be *GetResponse."
	}
	if s = ShouldBeValidSession(gr.Session); s != "" {
		return
	}
	return
}

func ShouldBeValidContextResponse(actual interface{}, expected ...interface{}) (s string) {
	if len(expected) != 1 {
		return fmt.Sprintf("This assertion requires exactly 1 comparison values (you provided %d).", len(expected))
	}

	cr, ok := actual.(*mnemosynerpc.ContextResponse)
	if !ok {
		return "The given value must be *ContextResponse."
	}
	if s = ShouldBeValidSession(cr.Session); s != "" {
		return
	}
	return
}

func ShouldBeValidSession(expected ...interface{}) (s string) {
	if len(expected) != 1 {
		return fmt.Sprintf("This assertion requires exactly 1 comparison values (you provided %d).", len(expected))
	}
	sr, ok := expected[0].(*mnemosynerpc.Session)
	if !ok {
		return "The given value must be *GetResponse."
	}
	if s = convey.ShouldNotBeNil(sr); s != "" {
		return
	}
	if s = convey.ShouldHaveLength(string(sr.AccessToken), 128); s != "" {
		return
	}
	return
}

func ShouldBeGRPCError(compare func(interface{}, ...interface{}) string) func(actual interface{}, expected ...interface{}) (s string) {
	return func(actual interface{}, expected ...interface{}) (s string) {
		if len(expected) != 2 {
			return fmt.Sprintf("This assertion requires exactly 2 comparison values (you provided %d).", len(expected))
		}

		e, ok := actual.(error)
		if !ok {
			return "The given value must implement error interface."
		}
		if s = compare(status.Convert(e).Message(), expected[1]); s != "" {
			return
		}
		if s = convey.ShouldEqual(status.Code(e), expected[0]); s != "" {
			return
		}
		return
	}
}

func ShouldBeValidToken(actual interface{}, expected ...interface{}) (s string) {
	if len(expected) != 0 {
		return fmt.Sprintf("This assertion requires exactly 0 comparison values (you provided %d).", len(expected))
	}

	if s = convey.ShouldNotBeNil(actual); s != "" {
		return
	}
	if s = convey.ShouldHaveSameTypeAs(actual, ""); s != "" {
		return
	}

	t, ok := actual.(string)
	if s = convey.ShouldBeTrue(ok); s != "" {
		return
	}
	if s = convey.ShouldNotBeEmpty(t); s != "" {
		return
	}
	return
}
