package echo_serv

import (
	"context"
	"echoserv/echo_pb"

	log "github.com/Golang-Tools/loggerhelper"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//Square req-res模式方法实现模板
func (s *Server) Echo(ctx context.Context, in *echo_pb.Message) (*echo_pb.Message, error) {
	log.Info("Echo get message", log.Dict{"in": in})
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Info("get header from client", log.Dict{"header": md})
	}
	header := metadata.Pairs("header-key", "val")
	grpc.SendHeader(ctx, header) //发送header
	// create and set trailer
	defer func() {
		trailer := metadata.Pairs("trailer-key", "val")
		grpc.SetTrailer(ctx, trailer) //设置trailer在函数执行完成后发送
	}()
	m := &echo_pb.Message{Message: in.Message}
	log.Info("Square send message", log.Dict{"result": m})
	return m, nil
}
