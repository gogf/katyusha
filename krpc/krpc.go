package krpc

type (
	krpcCtx    struct{}
	krpcClient struct{}
	krpcServer struct{}
)

const (
	configNodeNameGrpcServer = "grpcserver"
	configNodeNameHttpServer = "httpserver"
)

var (
	Ctx    = new(krpcCtx)  // Ctx manages the context feature.
	Client = &krpcClient{} // Client manages the client features.
	Server = &krpcServer{} // Server manages the server feature.
)
