package krpc

import (
	"context"
	"github.com/gogf/gf/util/gutil"
	"google.golang.org/grpc"
)

// UnaryRecover is the first interpreter that keep server not down from panics.
func (s *GrpcServer) UnaryRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	gutil.TryCatch(func() {
		res, err = handler(ctx, req)
	}, func(exception error) {
		err = exception
	})
	return
}
