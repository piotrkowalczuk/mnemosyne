package mnemosyned

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func errorInterceptor(logger *zap.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	{
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			res, err := handler(ctx, req)

			if err != nil && grpc.Code(err) != codes.OK {
				loggerBackground(ctx, logger).Error("request failure",
					zap.String("error", grpc.ErrorDesc(err)),
					zap.String("handler", info.FullMethod),
				)

				switch err {
				case errMissingSession, errMissingAccessToken, errMissingSubjectID, errSessionNotFound:
					return nil, err
				default:
					return nil, grpc.Errorf(grpc.Code(err), "mnemosyned: %s", grpc.ErrorDesc(err))
				}
			}

			loggerBackground(ctx, logger).Debug("request handled successfully", zap.String("handler", info.FullMethod))
			return res, err
		}
	}
}

func loggerBackground(ctx context.Context, logger *zap.Logger, fields ...zapcore.Field) *zap.Logger {
	l := logger.With(fields...)
	if md, ok := metadata.FromContext(ctx); ok {
		if rid, ok := md["request_id"]; ok && len(rid) >= 1 {
			l = l.With(zap.String("request_id", rid[0]))
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		l = l.With(zap.String("peer_address", p.Addr.String()))
	}
	return l
}
