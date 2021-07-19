package interceptor

import (
	"context"
	"time"

	log "github.com/Golang-Tools/loggerhelper"
	"google.golang.org/grpc"
)

// UnaryTimerIntercepto 一元拦截器
func UnaryTimerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	m, err := handler(ctx, req)
	end := time.Now()
	log.Info("rpc desc", log.Dict{"start": start.Format(time.RFC3339), "end": end.Format(time.RFC3339), "info": info, "err": err})
	return m, err
}