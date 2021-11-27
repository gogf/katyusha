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
	AppId      string // (necessary) Unique id for the service, only for service discovery.
	Address    string // (necessary) Service address, single one, usually IP:port, eg: 192.168.1.2:8000
	Deployment string // (optional)  Service deployment name, eg: dev, qa, staging, prod, etc.
	Group      string // (optional)  Service group, to indicate different service in the same environment with the same Name and AppId.
	Version    string // (optional)  Service version, eg: v1.0.0, v2.1.1, etc.
	Metadata   g.Map  // (optional)  Custom data for this service, which can be set using JSON by environment or command-line.
}

type discoveryEnvKey struct {
	PrefixRoot string
	KeepAlive  string
	AppId      string
	Address    string
	Version    string
	Deployment string
	Group      string
	Metadata   string
	Endpoints  string
}

type discoveryDefaultValue struct {
	PrefixRoot string
	KeepAlive  time.Duration
	Version    string
	Deployment string
	Group      string
	Scheme     string
}

const (
	configNodeNameDiscovery = "discovery"
	configNodeNameService   = "service"
)

var (
	EnvKey = discoveryEnvKey{
		PrefixRoot: "KA_PREFIX_ROOT",
		KeepAlive:  "KA_KEEPALIVE",
		AppId:      "KA_APP_ID",
		Address:    "KA_ADDRESS",
		Version:    "KA_VERSION",
		Deployment: "KA_DEPLOYMENT",
		Group:      "KA_GROUP",
		Metadata:   "KA_METADATA",
		Endpoints:  "KA_ENDPOINTS",
	}

	DefaultValue = discoveryDefaultValue{
		PrefixRoot: "/katyusha",
		KeepAlive:  10 * time.Second,
		Version:    "v0.0.0",
		Deployment: "default",
		Group:      "default",
		Scheme:     "katyusha",
	}
)
