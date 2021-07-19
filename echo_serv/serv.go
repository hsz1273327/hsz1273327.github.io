package echo_serv

import (
	"crypto/tls"
	"crypto/x509"
	"echoserv/echo_pb"
	"io/ioutil"
	"net"
	"os"
	"time"

	log "github.com/Golang-Tools/loggerhelper"

	grpc "google.golang.org/grpc"
	channelz "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"echoserv/echo_serv/interceptor"
)

//Server grpc的服务器结构体
//服务集成了如下特性:
//设置收发最大消息长度
//健康检测
//gzip做消息压缩
//接口反射
//channelz支持
//TLS支持
//keep alive 支持
type Server struct {
	App_Name    string `json:"app_name,omitempty" jsonschema:"description=服务名"`
	App_Version string `json:"app_version,omitempty" jsonschema:"description=服务版本"`
	Address     string `json:"address,omitempty" jsonschema:"description=服务的主机和端口"`
	Log_Level   string `json:"log_level,omitempty" jsonschema:"description=项目的log等级,enum=TRACE,enum=DEBUG,enum=INFO,enum=WARN,enum=ERROR"`

	// 性能设置
	Max_Recv_Msg_Size                           int  `json:"max_recv_msg_size,omitempty" jsonschema:"description=允许接收的最大消息长度"`
	Max_Send_Msg_Size                           int  `json:"max_send_msg_size,omitempty" jsonschema:"description=允许发送的最大消息长度"`
	Initial_Window_Size                         int  `json:"initial_window_size,omitempty" jsonschema:"description=基于Stream的滑动窗口大小"`
	Initial_Conn_Window_Size                    int  `json:"initial_conn_window_size,omitempty" jsonschema:"description=基于Connection的滑动窗口大小"`
	Max_Concurrent_Streams                      int  `json:"max_concurrent_streams,omitempty" jsonschema:"description=一个连接中最大并发Stream数"`
	Max_Connection_Idle                         int  `json:"max_connection_idle,omitempty" jsonschema:"description=客户端连接的最大空闲时长"`
	Max_Connection_Age                          int  `json:"max_connection_age,omitempty" jsonschema:"description=如果连接存活超过n则发送goaway"`
	Max_Connection_Age_Grace                    int  `json:"max_connection_age_grace,omitempty" jsonschema:"description=强制关闭连接之前允许等待的rpc在n秒内完成"`
	Keepalive_Time                              int  `json:"keepalive_time,omitempty" jsonschema:"description=空闲连接每隔n秒ping一次客户端已确保连接存活"`
	Keepalive_Timeout                           int  `json:"keepalive_timeout,omitempty" jsonschema:"description=ping时长超过n则认为连接已死"`
	Keepalive_Enforcement_Min_Time              int  `json:"keepalive_enforement_min_time,omitempty" jsonschema:"description=如果客户端超过每n秒ping一次则终止连接"`
	Keepalive_Enforcement_Permit_Without_Stream bool `json:"keepalive_enforement_permit_without_stream,omitempty" jsonschema:"description=即使没有活动流也允许ping"`

	//TLS设置
	// XDS_CREDS        bool   `json:"xds_creds,omitempty" jsonschema:"description=是否使用xDS APIs来接收TLS设置"`
	Server_Cert_Path string `json:"server_cert_path,omitempty" jsonschema:"description=使用TLS时服务端的证书位置"`
	Server_Key_Path  string `json:"server_key_path,omitempty" jsonschema:"description=使用TLS时服务端证书的私钥位置"`
	Ca_Cert_Path     string `json:"ca_cert_path,omitempty" jsonschema:"description=使用TLS时根整数位置"`
	Client_Crl_Path  string `json:"client_crl_path,omitempty" jsonschema:"description=客户端证书黑名单位置"`

	// 调试,目前admin接口不稳定,因此只使用channelz
	Use_Admin bool `json:"use_admin,omitempty" jsonschema:"description=是否使用grpc-admin方便调试"`

	opts          []grpc.ServerOption
	healthservice *health.Server
}

//Main 服务的入口函数
func (s *Server) Main() {
	// 初始化log
	log.Init(s.Log_Level, log.Dict{
		"app_name":    s.App_Name,
		"app_version": s.App_Version,
	})
	log.Info("获得参数", nil, log.Dict{"ServiceConfig": s}, nil)
	s.Run()
}

//PerformanceOpts 配置性能调优设置
func (s *Server) PerformanceOpts() {
	if s.opts == nil {
		s.opts = []grpc.ServerOption{}
	}

	if s.Max_Recv_Msg_Size != 0 {
		s.opts = append(s.opts, grpc.MaxRecvMsgSize(s.Max_Recv_Msg_Size))
	}
	if s.Max_Send_Msg_Size != 0 {
		s.opts = append(s.opts, grpc.MaxSendMsgSize(s.Max_Send_Msg_Size))
	}
	if s.Initial_Window_Size != 0 {
		s.opts = append(s.opts, grpc.InitialWindowSize(int32(s.Initial_Window_Size)))
	}
	if s.Initial_Conn_Window_Size != 0 {
		s.opts = append(s.opts, grpc.InitialConnWindowSize(int32(s.Initial_Conn_Window_Size)))
	}
	if s.Max_Concurrent_Streams != 0 {
		s.opts = append(s.opts, grpc.MaxConcurrentStreams(uint32(s.Max_Concurrent_Streams)))
	}
	if s.Max_Connection_Idle != 0 || s.Max_Connection_Age != 0 || s.Max_Connection_Age_Grace != 0 || s.Keepalive_Time != 0 || s.Keepalive_Timeout != 0 {
		kasp := keepalive.ServerParameters{
			MaxConnectionIdle:     time.Duration(s.Max_Connection_Idle) * time.Second,
			MaxConnectionAge:      time.Duration(s.Max_Connection_Age) * time.Second,
			MaxConnectionAgeGrace: time.Duration(s.Max_Connection_Age_Grace) * time.Second,
			Time:                  time.Duration(s.Keepalive_Time) * time.Second,
			Timeout:               time.Duration(s.Keepalive_Timeout) * time.Second,
		}
		s.opts = append(s.opts, grpc.KeepaliveParams(kasp))
	}

	if s.Keepalive_Enforcement_Min_Time != 0 || s.Keepalive_Enforcement_Permit_Without_Stream == true {
		kaep := keepalive.EnforcementPolicy{
			MinTime:             time.Duration(s.Keepalive_Enforcement_Min_Time) * time.Second,
			PermitWithoutStream: s.Keepalive_Enforcement_Permit_Without_Stream,
		}
		s.opts = append(s.opts, grpc.KeepaliveEnforcementPolicy(kaep))
	}
}

//TLSOpts 配置TLS设置
func (s *Server) TLSOpts() {
	if s.opts == nil {
		s.opts = []grpc.ServerOption{}
	}
	if s.Server_Cert_Path != "" && s.Server_Key_Path != "" {
		if s.Ca_Cert_Path != "" {
			cert, err := tls.LoadX509KeyPair(s.Server_Cert_Path, s.Server_Key_Path)
			if err != nil {
				log.Error("read serv pem file error:", log.Dict{"err": err.Error(), "Cert_path": s.Server_Cert_Path, "Key_Path": s.Server_Key_Path})
				os.Exit(2)
			}
			capool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(s.Ca_Cert_Path)
			if err != nil {
				log.Error("read ca pem file error:", log.Dict{"err": err.Error(), "path": s.Ca_Cert_Path})
				os.Exit(2)
			}
			capool.AppendCertsFromPEM(caCrt)
			tlsconf := &tls.Config{
				RootCAs:      capool,
				ClientAuth:   tls.RequireAndVerifyClientCert, // 检验客户端证书
				Certificates: []tls.Certificate{cert},
			}
			if s.Client_Crl_Path != "" {
				clipool := x509.NewCertPool()
				cliCrt, err := ioutil.ReadFile(s.Client_Crl_Path)
				if err != nil {
					log.Error("read pem file error:", log.Dict{"err": err.Error(), "path": s.Client_Crl_Path})
					os.Exit(2)
				}
				clipool.AppendCertsFromPEM(cliCrt)
				tlsconf.ClientCAs = clipool
			}
			creds := credentials.NewTLS(tlsconf)
			s.opts = append(s.opts, grpc.Creds(creds))
		} else {
			creds, err := credentials.NewServerTLSFromFile(s.Server_Cert_Path, s.Server_Key_Path)
			if err != nil {
				log.Error("Failed to Listen as a TLS Server", log.Dict{"error": err.Error()})
				os.Exit(2)
			}
			s.opts = append(s.opts, grpc.Creds(creds))
		}
	}
}

//RegistInterceptor 注册拦截器
func (s *Server) RegistInterceptor() {
	s.opts = append(s.opts,
		grpc.UnaryInterceptor(interceptor.UnaryTimerInterceptor),
		grpc.StreamInterceptor(interceptor.StreamTimerInterceptor),
	)
}

//RunServer 启动服务
func (s *Server) RunServer() {
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Error("Failed to Listen", log.Dict{"error": err.Error(), "address": s.Address})
		os.Exit(1)
	}
	s.PerformanceOpts()
	s.TLSOpts()
	s.RegistInterceptor()
	gs := grpc.NewServer(s.opts...)
	defer gs.Stop()
	// 注册健康检查
	s.healthservice = health.NewServer()
	healthpb.RegisterHealthServer(gs, s.healthservice)

	reflection.Register(gs)
	if s.Use_Admin {
		channelz.RegisterChannelzServiceToServer(gs)
	}
	// 注册服务
	echo_pb.RegisterECHOServer(gs, s)

	log.Info("Server Start", log.Dict{"address": s.Address})
	err = gs.Serve(lis)
	if err != nil {
		log.Error("Failed to Serve", log.Dict{"error": err})
		os.Exit(1)
	}
}

//Run 执行grpc服务
func (s *Server) Run() {
	s.RunServer()
}

var ServNode = Server{
	App_Name:    "echoserv",
	App_Version: "0.0.0",
	Address:     "0.0.0.0:5000",
	Log_Level:   "DEBUG",
}
