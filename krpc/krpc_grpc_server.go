package krpc

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gipv4"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/katyusha/discovery"
	"google.golang.org/grpc"
	"net"
	"sync"
)

// GrpcServer is the server for GRPC protocol.
type GrpcServer struct {
	Server    *grpc.Server
	config    GrpcServerConfig
	services  []*discovery.Service
	waitGroup sync.WaitGroup
}

// GrpcServerConfig is the configuration for server.
type GrpcServerConfig struct {
	Address string // Address for server listening.
}

// GrpcServerOption is alias for grpc.ServerOption.
type GrpcServerOption = grpc.ServerOption

// NewGrpcServer creates and returns a grpc server.
func NewGrpcServer(config GrpcServerConfig, option ...GrpcServerOption) *GrpcServer {
	if config.Address == "" {
		panic("server address cannot be empty")
	}
	if !gstr.Contains(config.Address, ":") {
		panic("invalid service address, should contain listening port")
	}
	server := &GrpcServer{
		Server: grpc.NewServer(option...),
		config: config,
	}
	return server
}

// Service binds service list to current server.
// Server will automatically register the service list after it starts.
func (s *GrpcServer) Service(services ...*discovery.Service) {
	var (
		serviceAddress string
		array          = gstr.Split(s.config.Address, ":")
	)
	if array[0] == "0.0.0.0" || array[0] == "" {
		intraIp, err := gipv4.GetIntranetIp()
		if err != nil {
			panic("retrieving intranet ip failed, please check your net card or manually assign the service address: " + err.Error())
		}
		serviceAddress = fmt.Sprintf(`%s:%s`, intraIp, array[1])
	} else {
		serviceAddress = s.config.Address
	}
	for _, service := range services {
		if service.Address == "" {
			service.Address = serviceAddress
		}
	}
	s.services = services
}

// Run starts the server in blocking way.
func (s *GrpcServer) Run() error {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}
	if len(s.services) == 0 {
		appId := gcmd.GetWithEnv(discovery.EnvKeyAppId).String()
		if appId != "" {
			// Automatically creating service if app id can be retrieved
			// from environment or command-line.
			s.Service(&discovery.Service{
				AppId: appId,
			})
		}
	}
	// Register service list after server starts.
	for _, service := range s.services {
		if err = discovery.Register(service); err != nil {
			return err
		}
	}
	g.Log().Printf("grpc server start listening on: %s", s.config.Address)
	return s.Server.Serve(listener)
}

// Start starts the server in no-blocking way.
func (s *GrpcServer) Start() {
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		if err := s.Run(); err != nil {
			panic(err)
		}
	}()
}

// Wait works with Start, which blocks current goroutine until the server stops.
func (s *GrpcServer) Wait() {
	s.waitGroup.Wait()
}

// Stop gracefully stops the server.
func (s *GrpcServer) Stop() {
	s.Server.GracefulStop()
}
