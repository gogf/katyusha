package discovery

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/text/gstr"
)

// Service definition.
type Service struct {
	AppId      string // (necessary) Unique id for the service, only for service discovery.
	Address    string // (necessary) Service address, usually IP:port.
	Deployment string // (optional)  Service deployment name, eg: dev, qa, staging, prod, etc.
	Group      string // (optional)  Service group, to indicate different service in the same environment with the same Name and AppId.
	Version    string // (optional)  Service version, eg: v1.0.0, v2.1.1, etc.
	Metadata   g.Map  // (optional)  Custom data for this service, which can be set using JSON by environment or command-line.
}

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
		err := json.Unmarshal(value, &service.Metadata)
		if err != nil {
			g.Log().Error(err)
		}
	}
	return service
}
