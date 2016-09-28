package cluster

import (
	"sort"

	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"google.golang.org/grpc"
)

// Cluster ...
type Cluster struct {
	listen  string
	buckets int
	nodes   []*Node
}

// New ...
func New(listen string, seeds ...string) (csr *Cluster, err error) {
	l := len(seeds) + 1
	nodes := make([]string, 0, l)
	nodes = append(nodes, listen)
	nodes = append(nodes, seeds...)
	sort.Strings(nodes)

	csr = &Cluster{
		nodes:   make([]*Node, 0, l),
		buckets: l,
		listen:  listen,
	}

	for _, addr := range nodes {
		csr.nodes = append(csr.nodes, &Node{
			Addr: addr,
		})
	}
	return csr, nil
}

// Connect ...
func (c *Cluster) Connect(opts ...grpc.DialOption) error {
	for _, n := range c.nodes {
		if n.Addr == c.listen {
			continue
		}

		conn, err := grpc.Dial(n.Addr, opts...)
		if err != nil {
			return err
		}

		n.Client = mnemosynerpc.NewSessionManagerClient(conn)
	}
	return nil
}

// Get if possible returns node that corresponds to given access token.
func (c *Cluster) Get(k int32) (*Node, bool) {
	if len(c.nodes) == 0 {
		return nil, false
	}
	if len(c.nodes)-1 < int(k) {
		return nil, false
	}
	return c.nodes[k], true
}

// Nodes ...
func (c *Cluster) Nodes() []*Node {
	return c.nodes
}

// Len returns number of nodes.
func (c *Cluster) Len() int {
	return c.buckets
}
