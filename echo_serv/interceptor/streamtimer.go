package interceptor

import (
	"fmt"
	"time"

	log "github.com/Golang-Tools/loggerhelper"
	"google.golang.org/grpc"
)

// wrappedStream 服务端流的包装
// SendMsg method call.
type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Info("Receive a message", log.Dict{"type": fmt.Sprintf("%T", m), "at": time.Now().Format(time.RFC3339)})
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Info("Send a message", log.Dict{"type": fmt.Sprintf("%T", m), "at": time.Now().Format(time.RFC3339)})
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// StreamInterceptor 流式拦截器
func StreamTimerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	err := handler(srv, newWrappedStream(ss))
	end := time.Now()
	log.Info("rpc desc", log.Dict{"start": start.Format(time.RFC3339), "end": end.Format(time.RFC3339), "info": info, "err": err})
	return err
}