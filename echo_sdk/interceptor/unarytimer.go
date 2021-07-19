package interceptor

import (
	"context"
	"time"

	log "github.com/Golang-Tools/loggerhelper"
	"google.golang.org/grpc"
)

// UnaryTimerIntercepto 一元拦截器
func UnaryTimerInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	log.Info("rpc desc", log.Dict{"start": start.Format(time.RFC3339), "end": end.Format(time.RFC3339), "method": method, "err": err})
	return err
}