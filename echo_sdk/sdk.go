package echo_sdk

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"echoserv/echo_pb"
	"echoserv/echo_sdk/interceptor"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/Golang-Tools/loggerhelper"
	jsoniter "github.com/json-iterator/go"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	resolver "google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	_ "google.golang.org/grpc/xds"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//SDKConfig 的客户端类型
type SDKConfig struct {
	Address      []string `json:"address" jsonschema:"description=连接服务的主机和端口"`
	Service_Name string   `json:"service_name,omitempty" jsonschema:"description=服务器域名"`
	App_Name     string   `json:"app_name,omitempty" jsonschema:"description=服务名"`
	App_Version  string   `json:"app_version,omitempty" jsonschema:"description=服务版本"`

	// 性能设置
	Initial_Window_Size                         int  `json:"initial_window_size,omitempty" jsonschema:"description=基于Stream的滑动窗口大小"`
	Initial_Conn_Window_Size                    int  `json:"initial_conn_window_size,omitempty" jsonschema:"description=基于Connection的滑动窗口大小"`
	Keepalive_Time                              int  `json:"keepalive_time,omitempty" jsonschema:"description=空闲连接每隔n秒ping一次客户端已确保连接存活"`
	Keepalive_Timeout                           int  `json:"keepalive_timeout,omitempty" jsonschema:"description=ping时长超过n则认为连接已死"`
	Keepalive_Enforcement_Permit_Without_Stream bool `json:"keepalive_enforement_permit_without_stream,omitempty" jsonschema:"description=是否当连接空闲时仍然发送PING帧监测"`
	Conn_With_Block                             bool `json:"conn_with_block,omitempty" jsonschema:"description=同步的连接建立"`
	Max_Recv_Msg_Size                           int  `json:"max_rec_msg_size,omitempty" jsonschema:"description=允许接收的最大消息长度"`
	Max_Send_Msg_Size                           int  `json:"max_send_msg_size,omitempty" jsonschema:"description=允许发送的最大消息长度"`

	//压缩设置,目前只支持gzip
	Compression string `json:"compression,omitempty" jsonschema:"description=使用哪种方式压缩发送的消息,enum=gzip"`

	// TLS设置
	Ca_Cert_Path     string `json:"ca_cert_path,omitempty" jsonschema:"description=如果要使用tls则需要指定根证书位置"`
	Client_Cert_Path string `json:"client_cert_path,omitempty" jsonschema:"description=客户端整数位置"`
	Client_Key_Path  string `json:"client_key_path,omitempty" jsonschema:"description=客户端证书对应的私钥位置"`

	// 请求超时设置
	Query_Timeout int `json:"query_timeout,omitempty" jsonschema:"description=请求服务的最大超时时间单位ms"`
}

//NewSDK 创建客户端对象
func (c *SDKConfig) NewSDK() *SDK {
	sdk := New()
	sdk.Init(c)
	return sdk
}

//SDK 的客户端类型
type SDK struct {
	*SDKConfig
	opts          []grpc.DialOption
	serviceconfig map[string]interface{}
	addr          string
}

//New 创建客户端对象
func New() *SDK {
	c := new(SDK)
	c.opts = make([]grpc.DialOption, 0, 10)
	return c
}

//InitCallOpts 初始化连接选项
func (c *SDK) InitCallOpts() {
	callopts := []grpc.CallOption{}
	if c.Max_Recv_Msg_Size != 0 {
		callopts = append(callopts, grpc.MaxCallRecvMsgSize(c.Max_Recv_Msg_Size))
	}
	if c.Max_Send_Msg_Size != 0 {
		callopts = append(callopts, grpc.MaxCallSendMsgSize(c.Max_Send_Msg_Size))
	}
	switch c.Compression {
	case "gzip":
		{
			callopts = append(callopts, grpc.UseCompressor(gzip.Name))
		}
	}

	if len(callopts) > 0 {
		c.opts = append(c.opts, grpc.WithDefaultCallOptions(callopts...))
	}
}

//InitOpts 初始化连接选项
func (c *SDK) InitOpts() error {
	if c.Ca_Cert_Path != "" {
		if c.Client_Cert_Path != "" && c.Client_Key_Path != "" {
			cert, err := tls.LoadX509KeyPair(c.Client_Cert_Path, c.Client_Key_Path)
			if err != nil {
				log.Error("read client pem file error:", log.Dict{"err": err.Error(), "Cert_path": c.Client_Cert_Path, "Key_Path": c.Client_Key_Path})
				return err
			}
			capool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(c.Ca_Cert_Path)
			if err != nil {
				log.Error("read ca pem file error:", log.Dict{"err": err.Error(), "path": c.Ca_Cert_Path})
				return err
			}
			capool.AppendCertsFromPEM(caCrt)
			tlsconf := &tls.Config{
				RootCAs:      capool,
				Certificates: []tls.Certificate{cert},
			}
			creds := credentials.NewTLS(tlsconf)
			c.opts = append(c.opts, grpc.WithTransportCredentials(creds))
		} else {
			creds, err := credentials.NewClientTLSFromFile(c.Ca_Cert_Path, "")
			if err != nil {
				log.Error("failed to load credentials", log.Dict{"err": err.Error()})
				return err
			}
			c.opts = append(c.opts, grpc.WithTransportCredentials(creds))
		}
	} else {
		c.opts = append(c.opts, grpc.WithInsecure())
	}
	if c.Keepalive_Time != 0 || c.Keepalive_Timeout != 0 || c.Keepalive_Enforcement_Permit_Without_Stream == true {
		kacp := keepalive.ClientParameters{
			Time:                time.Duration(c.Keepalive_Time) * time.Second,
			Timeout:             time.Duration(c.Keepalive_Timeout) * time.Second,
			PermitWithoutStream: c.Keepalive_Enforcement_Permit_Without_Stream, // send pings even without active streams
		}
		c.opts = append(c.opts, grpc.WithKeepaliveParams(kacp))
	}
	if c.Conn_With_Block == true {
		c.opts = append(c.opts, grpc.WithBlock())
	}
	if c.Initial_Window_Size != 0 {
		c.opts = append(c.opts, grpc.WithInitialWindowSize(int32(c.Initial_Window_Size)))
	}
	if c.Initial_Conn_Window_Size != 0 {
		c.opts = append(c.opts, grpc.WithInitialConnWindowSize(int32(c.Initial_Conn_Window_Size)))
	}
	return nil
}

//RegistInterceptor 注册拦截器
func (c *SDK) RegistInterceptor() {
	c.opts = append(c.opts,
		grpc.WithUnaryInterceptor(interceptor.UnaryTimerInterceptor),
		grpc.WithStreamInterceptor(interceptor.StreamTimerInterceptor))
}

//Init 初始化sdk客户端的连接信息
func (c *SDK) Init(conf *SDKConfig) error {
	c.SDKConfig = conf
	if conf.Address == nil {
		return errors.New("必须至少有一个地址")
	}
	switch len(conf.Address) {
	case 0:
		{
			return errors.New("必须至少有一个地址")
		}
	case 1:
		{
			c.initStandalone()
		}
	default:
		{
			c.initWithLocalBalance()
		}
	}
	err := c.InitOpts()
	if err != nil {
		return err
	}
	c.InitCallOpts()
	if c.serviceconfig != nil {
		serviceconfig, err := json.MarshalToString(c.serviceconfig)
		if err != nil {
			return err
		}
		c.opts = append(c.opts, grpc.WithDefaultServiceConfig(serviceconfig))
	}
	c.RegistInterceptor()
	return nil
}

//InitStandalone 初始化单机服务的连接配置
func (c *SDK) initStandalone() error {
	c.addr = c.Address[0]
	if strings.HasPrefix(c.addr, "dns:///") {
		if c.serviceconfig == nil {
			c.serviceconfig = map[string]interface{}{
				"loadBalancingPolicy": "round_robin",
			}
		} else {
			c.serviceconfig["loadBalancingPolicy"] = "round_robin"
		}
	}
	return nil
}

//InitWithLocalBalance 初始化本地负载均衡的连接配置
func (c *SDK) initWithLocalBalance() error {
	serverName := ""
	if c.App_Name != "" {
		if c.App_Version != "" {
			serverName = fmt.Sprintf("%s-%s", c.App_Name, strings.ReplaceAll(c.App_Version, ".", "_"))
		} else {
			serverName = c.App_Name
		}

	}
	if c.serviceconfig == nil {
		c.serviceconfig = map[string]interface{}{
			"loadBalancingPolicy": "round_robin",
			"healthCheckConfig":   map[string]string{"serviceName": c.Service_Name},
		}
	} else {
		c.serviceconfig["loadBalancingPolicy"] = "round_robin"
		c.serviceconfig["healthCheckConfig"] = map[string]string{"serviceName": serverName}
	}
	r := manual.NewBuilderWithScheme("localbalancer")
	addresses := []resolver.Address{}
	for _, addr := range c.Address {
		addresses = append(addresses, resolver.Address{Addr: addr})
	}
	r.InitialState(resolver.State{
		Addresses: addresses,
	})
	c.addr = fmt.Sprintf("%s:///%s", r.Scheme(), serverName)
	c.opts = append(c.opts, grpc.WithResolvers(r))
	return nil
}

//NewConn 建立一个新的连接
func (c *SDK) NewConn() (*Conn, error) {
	conn, err := newConn(c, c.addr, c.opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *SDK) NewCtx() (ctx context.Context, cancel context.CancelFunc) {
	if c.SDKConfig.Query_Timeout > 0 {
		timeout := time.Duration(c.SDKConfig.Query_Timeout) * time.Millisecond
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	return
}

//Conn 客户端类
type Conn struct {
	echo_pb.ECHOClient
	conn *grpc.ClientConn
	sdk  *SDK
	once bool
}

func newConn(sdk *SDK, addr string, opts ...grpc.DialOption) (*Conn, error) {
	c := new(Conn)
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.ECHOClient = echo_pb.NewECHOClient(conn)
	return c, nil
}

//Close 断开连接
func (c *Conn) Close() error {
	return c.conn.Close()
}

//DefaultSDK 默认的sdk客户端
var DefaultSDK = New()
