package server

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/katyusha/registry"
	"google.golang.org/grpc"
	"net"
	"sync"
)

// GrpcServer is the server for GRPC protocol.
type GrpcServer struct {
	Server    *grpc.Server
	config    GrpcServerConfig
	register  registry.Register
	waitGroup sync.WaitGroup
}

// GrpcServerConfig is the configuration for server.
type GrpcServerConfig struct {
	Addr string // Address for server listening.
}

// GrpcServerOption is alias for grpc.ServerOption.
type GrpcServerOption = grpc.ServerOption

// NewGrpcServer creates and returns a grpc server.
func NewGrpcServer(config GrpcServerConfig, option ...GrpcServerOption) *GrpcServer {
	server := &GrpcServer{
		Server: grpc.NewServer(option...),
		config: config,
	}
	return server
}

// Run starts the server in blocking way.
func (s *GrpcServer) Run() error {
	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return err
	}
	g.Log().Printf("grpc start listening on: %s", s.config.Addr)
	return s.Server.Serve(listener)
}

// Start starts the server in no-blocking way.
func (s *GrpcServer) Start() {
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		s.Run()
	}()
}

// Wait works with Start, which blocks current goroutine until the server stops.
func (s *GrpcServer) Wait() {
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		s.Run()
	}()
}

// Stop gracefully stops the server.
func (s *GrpcServer) Stop() {
	s.Server.GracefulStop()
}
