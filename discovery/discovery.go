package discovery

import (
	"time"
)

// Discovery interface for service.
type Discovery interface {
	Register(service *Service) error
	Unregister(service *Service) error
	Services() []*Service
	Close() error
}

const (
	EnvKeyPrefixRoot = "KA_PREFIX_ROOT"
	EnvKeyKeepAlive  = "KA_KEEPALIVE"
	EnvKeyAppId      = "KA_APP_ID"
	EnvKeyAddress    = "KA_ADDRESS"
	EnvKeyVersion    = "KA_VERSION"
	EnvKeyDeployment = "KA_DEPLOYMENT"
	EnvKeyGroup      = "KA_GROUP"
	EnvKeyMetaData   = "KA_METADATA"
	EnvKeyEndpoints  = "KA_ENDPOINTS"
)

var (
	DefaultDeployment   = "default"
	DefaultGroup        = "default"
	DefaultVersion      = "v0.0.0"
	DefaultKeepAliveTtl = 10 * time.Second
	DefaultPrefixRoot   = "/katyusha"
	DefaultScheme       = "katyusha"
)
