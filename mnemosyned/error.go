package mnemosyned

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne/internal/service/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func errorInterceptor(log *zap.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	{
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			res, err := handler(ctx, req)

			code := grpc.Code(err)
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
					zap.String("error", grpc.ErrorDesc(err)),
					logger.Ctx(ctx, info, codes.OK),
				)

				switch err {
				case errMissingSession, errMissingAccessToken, errMissingSubjectID, errSessionNotFound:
					return nil, err
				default:
					return nil, grpc.Errorf(grpc.Code(err), "mnemosyned: %s", grpc.ErrorDesc(err))
				}
			}

			loggerBackground(ctx, log).Debug("request handled successfully", logger.Ctx(ctx, info, codes.OK))
			return res, err
		}
	}
}

func loggerBackground(ctx context.Context, log *zap.Logger, fields ...zapcore.Field) *zap.Logger {
	l := log.With(fields...)
	if md, ok := metadata.FromContext(ctx); ok {
		if rid, ok := md["request_id"]; ok && len(rid) >= 1 {
			l = l.With(zap.String("request_id", rid[0]))
		}
	}
	return l
}
