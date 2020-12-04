package registry

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"time"
)

// Service definition.
type Service struct {
	Name       string // Service Name, for manually readable display.
	AppId      string // Unique id for the service, only for service discovery.
	Version    string // Service version, eg: v1.0.0, v2.1.1, etc.
	Address    string // Service address, usually IP:port .
	Deployment string // Service deployment name, eg: dev, qa, staging, prod, etc.
	Group      string // Service group, to indicate different service in the same environment with the same Name and AppId.
	Metadata   g.Map  // Custom data for this service.
}

// Register for service.
type Register interface {
	Register(service *Service) error
	Unregister(service *Service) error
	Close() error
}

const (
	DeploymentDefault   = "default"
	DeploymentDev       = "dev"
	DeploymentQA        = "qa"
	DeploymentStaging   = "staging"
	DeploymentProd      = "prod"
	DefaultGroup        = "default"
	DefaultVersion      = "default"
	DefaultRegistryDir  = "root"
	DefaultKeepAliveTtl = 10 * time.Second
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
