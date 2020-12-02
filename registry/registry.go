package registry

import (
	"github.com/gogf/gf/frame/g"
)

// Service definition.
type Service struct {
	Name        string // Service Name.
	AppId       string // Unique id for the service, it can be the same with Name.
	Version     string // Service version, eg: v1.0.0, v2.1.1, etc.
	Address     string // Service address, usually IP:port .
	Environment string // Service deployment environment, eg: dev, qa, staging, prod, etc.
	Group       string // Service group, to indicate different service in the same environment with the same Name and AppId.
	Metadata    g.Map  // Custom data that will be passed from service to service.
}

// Register for service.
type Register interface {
	Register(service *Service) error
	Unregister(service *Service) error
	Close()
}
