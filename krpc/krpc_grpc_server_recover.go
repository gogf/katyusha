package krpc

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// UnaryRecover is the first interceptor that keep server not down from panics.
func (s *GrpcServer) UnaryRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	gutil.TryCatch(func() {
		res, err = handler(ctx, req)
	}, func(exception error) {
		err = gerror.WrapCode(int(codes.Internal), err, "panic recovered")
	})
	return
}
