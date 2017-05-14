package cluster

import (
	"sort"

	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Cluster ...
type Cluster struct {
	listen  string
	buckets int
	nodes   []*Node
	logger  *zap.Logger
}

// Opts ...
type Opts struct {
	Listen string
	Seeds  []string
	Logger *zap.Logger
}

// New ...
func New(opts Opts) (csr *Cluster, err error) {
	var (
		nodes  []string
		exists bool
	)
	for _, seed := range opts.Seeds {
		if seed == opts.Listen {
			exists = true
		}
	}
	if !exists {
		nodes = append(nodes, opts.Listen)
	}
	nodes = append(nodes, opts.Seeds...)
	sort.Strings(nodes)

	csr = &Cluster{
		nodes:  make([]*Node, 0),
		listen: opts.Listen,
		logger: opts.Logger,
	}

	for _, addr := range nodes {
		if addr == "" {
			continue
		}
		csr.buckets++
		csr.nodes = append(csr.nodes, &Node{
			Addr: addr,
		})
	}
	return csr, nil
}

// Connect ...
func (c *Cluster) Connect(opts ...grpc.DialOption) error {
	for i, n := range c.nodes {
		if n.Addr == c.listen {
			continue
		}

		if c.logger != nil {
			c.logger.Debug("cluster node attempt to connect", zap.String("address", n.Addr), zap.Int("index", i))
		}

		conn, err := grpc.Dial(n.Addr, opts...)
		if err != nil {
			return err
		}

		if c.logger != nil {
			c.logger.Debug("cluster node connection success", zap.String("address", n.Addr), zap.Int("index", i))
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
