package mnemosyned

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/internal/service/logger"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

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

func errorInterceptor(log *zap.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	{
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			now := time.Now()

			if md, ok := metadata.FromIncomingContext(ctx); ok {
				ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
					mnemosyne.AccessTokenMetadataKey: md[mnemosyne.AccessTokenMetadataKey],
					"request_id":                     md["request_id"],
				})
			}

			res, err := handler(ctx, req)

			code := status.Code(err)
			if err != nil && code != codes.OK {
				if code == codes.Unknown {
					switch err {
					case sql.ErrNoRows:
						code = codes.NotFound
					default:
						if pqerr, ok := err.(*pq.Error); ok {
							switch pqerr.Code {
							case pq.ErrorCode("57014"):
								code = codes.Canceled
							}
						} else {
							code = codes.Internal
						}
					}
				}
				loggerBackground(ctx, log).Error("request failure",
					zap.String("error", status.Convert(err).Message()),
					logger.Ctx(ctx, info, code),
				)

				switch err {
				case errMissingAccessToken, errMissingSession, errMissingSubjectID:
					return nil, err
				case storage.ErrSessionNotFound:
					return nil, status.Errorf(codes.NotFound, "mnemosyned: %s", err.Error())
				case storage.ErrMissingAccessToken, storage.ErrMissingSession, storage.ErrMissingSubjectID:
					return nil, status.Errorf(codes.InvalidArgument, "mnemosyned: %s", err.Error())
				default:
					return nil, status.Errorf(status.Code(err), "mnemosyned: %s", status.Convert(err).Message())
				}
			}

			loggerBackground(ctx, log).Debug("request handled successfully",
				logger.Ctx(ctx, info, codes.OK),
				zap.Duration("elapsed", time.Since(now)),
			)
			return res, err
		}
	}
}

func loggerBackground(ctx context.Context, log *zap.Logger, fields ...zapcore.Field) *zap.Logger {
	l := log.With(fields...)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if rid, ok := md["request_id"]; ok && len(rid) >= 1 {
			l = l.With(zap.String("request_id", rid[0]))
		}
	}
	return l
}
