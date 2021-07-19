package echo_sdk

import (
	"echoserv/echo_pb"
	"fmt"
	"os"
	"time"

	log "github.com/Golang-Tools/loggerhelper"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//Main 执行测试
func (c *SDKConfig) Main() {
	fmt.Println(c)
	sdk := c.NewSDK()
	Conn, err := sdk.NewConn()
	if err != nil {
		log.Error("new comm get err", log.Dict{"err": err.Error()})
		os.Exit(1)
	}
	defer Conn.Close()
	//req-res
	ctx, cancel := sdk.NewCtx()
	defer cancel()
	md := metadata.Pairs("a", "1", "b", "2")
	ctx_meta := metadata.NewOutgoingContext(ctx, md)
	var header, trailer metadata.MD
	count := 20
	for t := 0; t < count; t++ {
		req, err := Conn.Echo(ctx_meta, &echo_pb.Message{Message: "hello"},
			grpc.Header(&header),   // will retrieve header
			grpc.Trailer(&trailer), // will retrieve trailer
		)
		if err != nil {
			log.Error("Square get error", log.Dict{"err": err.Error()})
			os.Exit(1)
		}
		log.Info("Square get result", log.Dict{"header": header, "req": req, "trailer": trailer})
		time.Sleep(100 * time.Millisecond)
	}
}

var TestNode = SDKConfig{
	Address: []string{"localhost:5000"},
}
