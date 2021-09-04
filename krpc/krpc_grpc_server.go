// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package krpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gipv4"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/katyusha/discovery"
	"google.golang.org/grpc"
)

// GrpcServer is the server for GRPC protocol.
type GrpcServer struct {
	Server    *grpc.Server
	Logger    *glog.Logger
	config    *GrpcServerConfig
	services  []*discovery.Service
	waitGroup sync.WaitGroup
}

// NewGrpcServer creates and returns a grpc server.
func (s krpcServer) NewGrpcServer(conf ...*GrpcServerConfig) *GrpcServer {
	var config *GrpcServerConfig
	if len(conf) > 0 {
		config = conf[0]
	} else {
		config = s.NewGrpcServerConfig()
	}
	if config.Address == "" {
		randomPort := s.randomPort()
		if randomPort == randomPortNotAvailable {
			g.Log().Fatal("server address is empty and random port retrieving failed")
		}
		config.Address = fmt.Sprintf(`:%d`, randomPort)
	}
	if !gstr.Contains(config.Address, ":") {
		g.Log().Fatal("invalid service address, should contain listening port")
	}
	if config.Logger == nil {
		config.Logger = glog.New()
	}
	grpcServer := &GrpcServer{
		Logger: config.Logger,
		config: config,
	}
	grpcServer.config.Options = append([]grpc.ServerOption{
		s.ChainUnary(
			grpcServer.UnaryLogger,
			s.UnaryError,
			s.UnaryRecover,
		),
	}, grpcServer.config.Options...)
	grpcServer.Server = grpc.NewServer(grpcServer.config.Options...)
	return grpcServer
}

// randomPort returns a random port that is not used by other processes.
func (s krpcServer) randomPort() int {
	intranetIp, err := gipv4.GetIntranetIp()
	if err != nil {
		panic(err)
	}
	for i := randomPortMin; i <= randomPortMax; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf(`%s:%d`, intranetIp, i))
		if err != nil {
			return i
		} else {
			conn.Close()
		}
	}
	return randomPortNotAvailable
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
			s.Logger.Fatal("retrieving intranet ip failed, please check your net card or manually assign the service address: " + err.Error())
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
	s.services = append(s.services, services...)
}

// Run starts the server in blocking way.
func (s *GrpcServer) Run() {
	if err := discovery.InitDiscoveryFromConfig(); err != nil {
		s.Logger.Fatal(err)
	}
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		s.Logger.Fatal(err)
	}
	if len(s.services) == 0 {
		appId := gcmd.GetOptWithEnv(discovery.EnvKey.AppId).String()
		if appId != "" {
			// Automatically creating service if app id can be retrieved
			// from environment or command-line.
			s.Service(&discovery.Service{
				AppId: appId,
			})
		}
		// Check any application identities bound with server.
		if len(s.config.AppId) > 0 {
			s.Service(&discovery.Service{
				AppId: s.config.AppId,
			})
		}
	}
	// Start listening.
	go func() {
		if err := s.Server.Serve(listener); err != nil {
			s.Logger.Fatal(err)
		}
	}()

	// Register service list after server starts.
	for _, service := range s.services {
		if err = discovery.Register(service); err != nil {
			s.Logger.Fatal(err)
		}
	}

	s.Logger.Printf("grpc server start listening on: %s, pid: %d", s.config.Address, gproc.Pid())

	// Signal listening and handling for gracefully shutdown.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)
	for {
		sig := <-sigChan
		switch sig {
		case
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGABRT:
			s.Logger.Printf("signal received: %s, gracefully shutting down", sig.String())
			for _, service := range s.services {
				_ = discovery.Unregister(service)
			}
			time.Sleep(time.Second)
			s.Stop()
			return
		}
	}
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
	s.waitGroup.Wait()
}

// Stop gracefully stops the server.
func (s *GrpcServer) Stop() {
	s.Server.GracefulStop()
}
