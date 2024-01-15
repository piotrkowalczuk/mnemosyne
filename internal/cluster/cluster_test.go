package cluster_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"go.uber.org/zap"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/mnemosyned"
)

var testPostgresAddress string

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

func TestNew(t *testing.T) {
	create := func(t *testing.T, args ...string) *cluster.Cluster {
		c, err := cluster.New(cluster.Opts{
			Listen: args[0], Seeds: args[1:],
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		return c
	}
	max := 100
	args := []string{"172.17.0.1", "172.17.0.2", "172.17.0.3", "127.0.0.1", "10.0.0.1", "8.8.8.8"}
	cs := make([]*cluster.Cluster, 0, max)
	for i := 0; i < max; i++ {
		cs = append(cs, create(t, args...))
	}

	for i, c := range cs {
		if i == 0 {
			continue
		}
		for j, n := range c.Nodes() {
			nodes := cs[i-1].Nodes()
			if nodes[j].Addr != n.Addr {
				t.Errorf("node address mismatch: %d %d: %s %s", i, j, nodes[j].Addr, n.Addr)
			}
		}
		if c.Len() != len(args) {
			t.Errorf("wrong number of nodes, expected %d but got %d", len(args), c.Len())
		}
	}
}

func TestCluster_Get_empty(t *testing.T) {
	c, err := cluster.New(cluster.Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, ok := c.Get(0)
	if ok {
		t.Error("expected nothing")
	}
}

func TestCluster_Get_beyond(t *testing.T) {
	c, err := cluster.New(cluster.Opts{
		Listen: "172.17.0.1",
		Seeds:  []string{"172.17.0.2", "172.17.0.3", "10.10.0.1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, ok := c.Get(5)
	if ok {
		t.Error("expected nothing")
	}
}

func TestCluster_Get(t *testing.T) {
	listen := "172.17.0.1"
	seeds := []string{"172.17.0.2", "172.17.0.3", "10.10.0.1"}
	c, err := cluster.New(cluster.Opts{
		Listen: listen,
		Seeds:  seeds,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	nodes := append(seeds, listen)
	sort.Strings(nodes)
	for k, addr := range nodes {
		got, ok := c.Get(int32(k))
		if !ok {
			t.Errorf("node not found: %s", addr)
			continue
		}
		if strings.HasPrefix(addr, "10") {
			continue
		}
		if got.Addr != addr {
			t.Errorf("address mismatch, expected %s but got %s", addr, got.Addr)
		} else {
			t.Logf("node under key %d and address %s passed", k, addr)
		}
	}
}

func TestCluster_GetOther_nil(t *testing.T) {
	var c *cluster.Cluster
	_, ok := c.GetOther("smth")
	if ok {
		t.Error("expected false")
	}
}

func TestCluster_GetOther_one(t *testing.T) {
	c, err := cluster.New(cluster.Opts{
		Listen: "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, ok := c.GetOther("smth")
	if ok {
		t.Error("expected false")
	}
}

func TestCluster_GetOther(t *testing.T) {
	c, cancel := testCluster(t)
	defer cancel()

	if err := c.Connect(context.TODO(), grpc.WithInsecure(), grpc.WithBlock()); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	for k := range c.Nodes() {
		got, ok := c.GetOther(fmt.Sprintf("access-token-%d", k))
		if ok {
			if got.Client == nil {
				t.Error("should not return node without established connection")
			}
			if got.Addr == c.Listen() {
				t.Error("should current node")
			}
		}
	}
}

func TestCluster_Connect(t *testing.T) {
	c, cancel := testCluster(t)
	defer cancel()
	if err := c.Connect(context.TODO(), grpc.WithInsecure(), grpc.WithBlock()); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func testCluster(t *testing.T) (*cluster.Cluster, func()) {
	m1, m1c := mnemosyned.TestDaemon(t, mnemosyned.TestDaemonOpts{
		StoragePostgresAddress: testPostgresAddress,
	})
	m2, m2c := mnemosyned.TestDaemon(t, mnemosyned.TestDaemonOpts{
		StoragePostgresAddress: testPostgresAddress,
	})
	m3, m3c := mnemosyned.TestDaemon(t, mnemosyned.TestDaemonOpts{
		StoragePostgresAddress: testPostgresAddress,
	})

	c, err := cluster.New(cluster.Opts{
		Listen: m1.String(),
		Seeds:  []string{m1.String(), m2.String(), m3.String()},
		Logger: zap.L(),
	})
	if err != nil {
		m1c.Close()
		m2c.Close()
		m3c.Close()
		t.Fatalf("unexpected error: %s", err.Error())
	}
	return c, func() {
		m3c.Close()
		m2c.Close()
		m1c.Close()
	}
}
