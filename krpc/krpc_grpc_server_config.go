package krpc

import (
	"github.com/gogf/gf/os/glog"
	"google.golang.org/grpc"
)

// GrpcServerConfig is the configuration for server.
type GrpcServerConfig struct {
	Address          string              // Address for server listening.
	Logger           *glog.Logger        // Logger for server.
	LogPath          string              // LogPath specifies the directory for storing logging files.
	LogStdout        bool                // LogStdout specifies whether printing logging content to stdout.
	ErrorLogEnabled  bool                // ErrorLogEnabled enables error logging content to files.
	ErrorLogPattern  string              // ErrorLogPattern specifies the error log file pattern like: error-{Ymd}.log
	AccessLogEnabled bool                // AccessLogEnabled enables access logging content to files.
	AccessLogPattern string              // AccessLogPattern specifies the error log file pattern like: access-{Ymd}.log
	Options          []grpc.ServerOption // Server options.
}

// NewGrpcServerConfig creates and returns a ServerConfig object with default configurations.
// Note that, do not define this default configuration to local package variable, as there are
// some pointer attributes that may be shared in different servers.
func NewGrpcServerConfig() *GrpcServerConfig {
	return &GrpcServerConfig{
		Address:          ":8000",
		Logger:           glog.New(),
		LogStdout:        true,
		ErrorLogEnabled:  true,
		ErrorLogPattern:  "error-{Ymd}.log",
		AccessLogEnabled: false,
		AccessLogPattern: "access-{Ymd}.log",
	}
}
