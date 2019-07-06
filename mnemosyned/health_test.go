package mnemosyned

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
)

func TestReadinessHandler_ServeHTTP(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	var (
		res *http.Response
		pay []byte
		err error
	)
	srv := httptest.NewServer(&readinessHandler{
		postgres: s.db,
		logger:   s.logger,
		cluster:  &cluster.Cluster{},
	})
	defer srv.Close()

	if res, err = http.Get(srv.URL); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("wrong status code, expected %d but got %d", http.StatusOK, res.StatusCode)
	}
	if pay, err = ioutil.ReadAll(res.Body); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	expPayload := `{"version":"","probes":{"postgres":"SERVING","cluster":{}}}
`
	if string(pay) != expPayload {
		t.Errorf("wrong payload, expected `%s` but got `%s`", expPayload, string(pay))
	}

	s.teardown(t)

	if res, err = http.Get(srv.URL); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	//if res.StatusCode != http.StatusServiceUnavailable {
	//	t.Fatalf("wrong status code, expected %d but got %d", http.StatusServiceUnavailable, res.StatusCode)
	//}
	if _, err = ioutil.ReadAll(res.Body); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	//if string(pay) != "postgres ping failure\n" {
	//	t.Errorf("wrong payload, expected '%s' but got '%s'", "postgres ping failure\n", string(pay))
	//}
}
