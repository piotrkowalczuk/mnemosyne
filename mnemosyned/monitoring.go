package mnemosyned

// asdasd
import (
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	monitoringRPCLabels = []string{
		"handler",
		"code",
	}
	monitoringPostgresLabels = []string{
		"query",
	}
)

type monitoring struct {
	enabled  bool
	cleanup  monitoringCleanup
	rpc      monitoringRPC
	postgres monitoringPostgres
	cache    monitoringCache
}

type monitoringCleanup struct {
	enabled bool
	errors  prometheus.Counter
}

type monitoringRPC struct {
	enabled  bool
	duration *prometheus.SummaryVec
	requests *prometheus.CounterVec
	errors   *prometheus.CounterVec
}

type monitoringPostgres struct {
	enabled bool
	queries *prometheus.CounterVec
	errors  *prometheus.CounterVec
}

type monitoringCache struct {
	enabled bool
	hits    prometheus.Counter
	misses  prometheus.Counter
	refresh prometheus.Counter
}

func unaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		wrap := func(current grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
			return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return current(currentCtx, currentReq, info, next)
			}
		}
		chain := handler
		for _, i := range interceptors {
			chain = wrap(i, chain)
		}
		return chain(ctx, req)
	}
}
