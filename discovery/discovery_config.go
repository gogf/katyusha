// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package discovery

import (
	"encoding/json"

	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
)

var (
	// initializedFromConfig is used for initialization for discovery.
	initializedFromConfig = gtype.NewBool()
)

// InitDiscoveryFromConfig automatically checks and initializes discovery feature
// from configuration.
func InitDiscoveryFromConfig() error {
	if !initializedFromConfig.Cas(false, true) {
		return nil
	}
	if !g.Cfg().Available() {
		return nil
	}
	// Configuration: discovery
	configDiscovery := g.Cfg().GetVar(configNodeNameDiscovery)
	if !configDiscovery.IsNil() {
		var (
			config  *Config
			service *Service
		)
		// Discovery.
		if err := configDiscovery.Struct(&config); err != nil {
			return err
		}
		if err := discoveryConfigToEnvironment(config); err != nil {
			return err
		}

		// Service.
		if err := configDiscovery.Struct(&service); err != nil {
			return err
		}
		if err := serviceConfigToEnvironment(service); err != nil {
			return err
		}
	}
	// Configuration: service
	configService := g.Cfg().GetVar(configNodeNameService)
	if !configService.IsNil() {
		if configService.IsSlice() {
			var (
				services []*Service
			)
			if err := configService.Structs(&services); err != nil {
				return err
			}
			for _, service := range services {
				if err := serviceConfigToEnvironment(service); err != nil {
					return err
				}
			}
		} else {
			var (
				service *Service
			)
			if err := configService.Struct(&service); err != nil {
				return err
			}
			if err := serviceConfigToEnvironment(service); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetConfig sets the discovery configuration using Config.
func SetConfig(config *Config) error {
	if err := discoveryConfigToEnvironment(config); err != nil {
		return err
	}
	return nil
}

// discoveryConfigToEnvironment sets the discovery environment value with Config object.
func discoveryConfigToEnvironment(config *Config) error {
	if config == nil {
		return nil
	}
	if len(config.Endpoints) > 0 {
		if err := genv.Set(EnvKey.Endpoints, gstr.Join(config.Endpoints, ",")); err != nil {
			return err
		}
	}
	if config.KeepAlive > 0 {
		if err := genv.Set(EnvKey.KeepAlive, config.KeepAlive.String()); err != nil {
			return err
		}
	}
	if config.PrefixRoot != "" {
		if err := genv.Set(EnvKey.PrefixRoot, config.PrefixRoot); err != nil {
			return err
		}
	}
	return nil
}

// serviceConfigToEnvironment sets the service environment value with Service object.
func serviceConfigToEnvironment(service *Service) error {
	if service == nil {
		return nil
	}
	if service.AppId != "" {
		if err := genv.Set(EnvKey.AppId, service.AppId); err != nil {
			return err
		}
	}
	if service.Address != "" {
		if err := genv.Set(EnvKey.Address, service.Address); err != nil {
			return err
		}
	}
	if service.Version != "" {
		if err := genv.Set(EnvKey.Version, service.Version); err != nil {
			return err
		}
	}
	if service.Group != "" {
		if err := genv.Set(EnvKey.Group, service.Group); err != nil {
			return err
		}
	}
	if service.Deployment != "" {
		if err := genv.Set(EnvKey.Deployment, service.Deployment); err != nil {
			return err
		}
	}
	if len(service.Metadata) > 0 {
		b, _ := json.Marshal(service.Metadata)
		if err := genv.Set(EnvKey.Metadata, string(b)); err != nil {
			return err
		}
	}
	return nil
}
