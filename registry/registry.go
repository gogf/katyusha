package registry

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	etcd3 "go.etcd.io/etcd/client/v3"
)

// Service definition.
type Service struct {
	Name        string // Service Name.
	AppId       string // Unique id for the service, it can be the same with Name.
	Version     string // Service version, eg: v1.0.0, v2.1.1, etc.
	Address     string // Service address, usually IP:port .
	Deployment  string // Service deployment name, eg: dev, qa, staging, prod, etc.
	Group       string // Service group, to indicate different service in the same environment with the same Name and AppId.
	Metadata    g.Map  // Custom data for this service.
	etcdGrantId etcd3.LeaseID
}

// Register for service.
type Register interface {
	Register(service *Service) error
	Unregister(service *Service) error
	Close() error
}

const (
	DeploymentDefault  = "default"
	DeploymentDev      = "dev"
	DeploymentQA       = "qa"
	DeploymentStaging  = "staging"
	DeploymentProd     = "prod"
	DefaultGroup       = "default"
	DefaultVersion     = "default"
	DefaultRegistryDir = "root"
)

// RegisterKey formats the service information with `registryDir` and returns the key string
// for registering.
func (s *Service) RegisterKey(registryDir string) string {
	return gstr.Join([]string{
		registryDir,
		s.Deployment,
		s.Group,
		s.AppId,
		s.Version,
	}, "/")
}
