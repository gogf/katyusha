package krpc

var (
	Ctx    = new(krpcCtx)  // Ctx manages the context feature.
	Client = &krpcClient{} // Client manages the client features.
	Server = &krpcServer{} // Server manages the server feature.
)

type (
	krpcCtx    struct{}
	krpcClient struct{}
	krpcServer struct{}
)
