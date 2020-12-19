package krpc

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"google.golang.org/grpc"
)

// GrpcServerConfig is the configuration for server.
type GrpcServerConfig struct {
	Address          string              // (necessary) Address for server listening.
	AppId            []string            // (optional)  Multiple application identities bound with this server.
	Logger           *glog.Logger        // (optional)  Logger for server.
	LogPath          string              // (optional)  LogPath specifies the directory for storing logging files.
	LogStdout        bool                // (optional)  LogStdout specifies whether printing logging content to stdout.
	ErrorLogEnabled  bool                // (optional)  ErrorLogEnabled enables error logging content to files.
	ErrorLogPattern  string              // (optional)  ErrorLogPattern specifies the error log file pattern like: error-{Ymd}.log
	AccessLogEnabled bool                // (optional)  AccessLogEnabled enables access logging content to files.
	AccessLogPattern string              // (optional)  AccessLogPattern specifies the error log file pattern like: access-{Ymd}.log
	Options          []grpc.ServerOption // (optional)  Server options.
}

// NewGrpcServerConfig creates and returns a ServerConfig object with default configurations.
// Note that, do not define this default configuration to local package variable, as there are
// some pointer attributes that may be shared in different servers.
func (s *krpcServer) NewGrpcServerConfig() *GrpcServerConfig {
	config := &GrpcServerConfig{
		Address:          ":8000",
		Logger:           glog.New(),
		LogStdout:        true,
		ErrorLogEnabled:  true,
		ErrorLogPattern:  "error-{Ymd}.log",
		AccessLogEnabled: false,
		AccessLogPattern: "access-{Ymd}.log",
	}
	// Reading configuration file and updating the configured keys.
	if g.Cfg().Available() {
		err := g.Cfg().GetVar(configNodeNameGrpcServer).Struct(&config)
		if err != nil {
			g.Log().Error(err)
		}
	}
	return config
}

// SetWithMap changes current configuration with map.
// This is commonly used for changing several configurations of current object.
func (c *GrpcServerConfig) SetWithMap(m g.Map) error {
	return gconv.Struct(m, c)
}

// MustSetWithMap acts as SetWithMap but panics if error occurs.
func (c *GrpcServerConfig) MustSetWithMap(m g.Map) {
	err := c.SetWithMap(m)
	if err != nil {
		panic(err)
	}
}
