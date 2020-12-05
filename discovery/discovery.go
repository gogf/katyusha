package discovery

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/text/gstr"
	"time"
)

// Service definition.
type Service struct {
	AppId      string // (necessary) Unique id for the service, only for service discovery.
	Address    string // (necessary) Service address, usually IP:port.
	Deployment string // (optional)  Service deployment name, eg: dev, qa, staging, prod, etc.
	Group      string // (optional)  Service group, to indicate different service in the same environment with the same Name and AppId.
	Version    string // (optional)  Service version, eg: v1.0.0, v2.1.1, etc.
	Metadata   g.Map  // (optional)  Custom data for this service.
}

// Registry interface for service.
type Registry interface {
	Register(service *Service) error
	Unregister(service *Service) error
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

// RegisterKey formats the service information with `registryDir` and returns the key string
// for registering.
func (s *Service) RegisterKey() string {
	return gstr.Join([]string{
		gcmd.GetWithEnv(EnvKeyPrefixRoot, DefaultPrefixRoot).String(),
		s.Deployment,
		s.Group,
		s.AppId,
		s.Version,
		s.Address,
	}, "/")
}

// newServiceFromKeyValue creates and returns service from `key` and `value`.
func newServiceFromKeyValue(key, value []byte) *Service {
	array := gstr.SplitAndTrim(string(key), "/")
	if len(array) < 6 {
		return nil
	}
	service := &Service{
		Deployment: array[1],
		Group:      array[2],
		AppId:      array[3],
		Version:    array[4],
		Address:    array[5],
		Metadata:   make(g.Map),
	}
	if len(value) > 0 {
		json.Unmarshal(value, &service.Metadata)
	}
	return service
}
