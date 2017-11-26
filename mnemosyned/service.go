package mnemosyned

import (
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"go.uber.org/zap"
)

func initCluster(l *zap.Logger, addr string, seeds ...string) (*cluster.Cluster, error) {
	csr, err := cluster.New(cluster.Opts{
		Listen: addr,
		Seeds:  seeds,
		Logger: l,
	})
	if err != nil {
		return nil, err
	}

	l.Debug("cluster initialized", zap.Strings("seeds", seeds), zap.String("listen", addr))
	return csr, nil
}
