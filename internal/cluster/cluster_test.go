package cluster_test

import (
	"testing"

	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/jump"
)

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

	for _, n := range c.Nodes() {
		hs := jump.HashString(n.Addr, len(seeds)+1)
		_, ok := c.Get(hs)
		if !ok {
			t.Errorf("node not found: %s", n.Addr)
			continue
		}
		//if got.Addr != n.Addr {
		//	t.Errorf("address mismatch, expected %s but got %s", n.Addr, got.Addr)
		//}
		//t.Logf("node under key %d and address %s passed", k, n.Addr)
	}
}
