package interceptor

import (
	"context"
	"fmt"
	"time"

	log "github.com/Golang-Tools/loggerhelper"
	"google.golang.org/grpc"
)

// wrappedStream  wraps around the embedded grpc.ClientStream, and intercepts the RecvMsg and
// SendMsg method call.
type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Info("Receive a message", log.Dict{"type": fmt.Sprintf("%T", m), "at": time.Now().Format(time.RFC3339)})
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Info("Send a message", log.Dict{"type": fmt.Sprintf("%T", m), "at": time.Now().Format(time.RFC3339)})
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

// StreamInterceptor 流式拦截器
func StreamTimerInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStream(s), nil
}