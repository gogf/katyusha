package krpc

import (
	"github.com/gogf/katyusha/krpc/internal/grpcctx"
)

type (
	krpcClient struct{}
	krpcServer struct{}
)

const (
	configNodeNameGrpcServer = "grpcserver"
	configNodeNameHttpServer = "httpserver"
)

var (
	Ctx    = grpcctx.Ctx   // Ctx manages the context feature.
	Client = &krpcClient{} // Client manages the client features.
	Server = &krpcServer{} // Server manages the server feature.
)
