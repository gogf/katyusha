// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package discovery

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Discovery interface for service.
type Discovery interface {
	Register(service *Service) error
	Unregister(service *Service) error
	Services() ([]*Service, error)
	Close() error
}

// Config is the configuration definition for discovery.
type Config struct {
	Endpoints  []string      // (necessary) The discovery server endpoints.
	PrefixRoot string        // (optional) Prefix string for discovery.
	KeepAlive  time.Duration // (optional) Keepalive duration for watcher.
}

// Service definition.
type Service struct {
	AppID      string // (necessary) Unique id for the service, only for service discovery.
	Address    string // (necessary) Service address, single one, usually IP:port, eg: 192.168.1.2:8000
	Deployment string // (optional)  Service deployment name, eg: dev, qa, staging, prod, etc.
	Group      string // (optional)  Service group, to indicate different service in the same environment with the same Name and AppID.
	Version    string // (optional)  Service version, eg: v1.0.0, v2.1.1, etc.
	Metadata   g.Map  // (optional)  Custom data for this service, which can be set using JSON by environment or command-line.
}

const (
	EnvPrefixRoot = "KA_PREFIX_ROOT"
	EnvKeepAlive  = "KA_KEEPALIVE"
	EnvAppID      = "KA_APP_ID"
	EnvAddress    = "KA_ADDRESS"
	EnvVersion    = "KA_VERSION"
	EnvDeployment = "KA_DEPLOYMENT"
	EnvGroup      = "KA_GROUP"
	EnvMetadata   = "KA_METADATA"
	EnvEndpoints  = "KA_ENDPOINTS"
)

const (
	DefaultPrefixRoot = "/katyusha"
	DefaultKeepAlive  = 10 * time.Second
	DefaultVersion    = "v0.0.0"
	DefaultDeployment = "default"
	DefaultGroup      = "default"
	DefaultScheme     = "katyusha"
)

const (
	configNodeNameDiscovery = "discovery"
	configNodeNameService   = "service"
)
