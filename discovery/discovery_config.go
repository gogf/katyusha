// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package discovery

import (
	"context"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	// initializedFromConfig is used for initialization of discovery.
	initializedFromConfig = gtype.NewBool()
)

// SetConfig sets the discovery configuration using Config.
func SetConfig(config *Config) error {
	if err := setDiscoveryConfigToEnvironment(config); err != nil {
		return err
	}
	return nil
}

// InitDiscoveryFromConfig automatically checks and initializes discovery feature
// from configuration.
func InitDiscoveryFromConfig() error {
	if !initializedFromConfig.Cas(false, true) {
		return nil
	}
	var (
		config  *Config
		service *Service
		ctx     = context.TODO()
	)
	// Configuration: discovery
	configDiscovery, err := g.Cfg().Get(ctx, configNodeNameDiscovery)
	if err != nil {
		return err
	}
	if !configDiscovery.IsNil() {
		// Discovery.
		if err = configDiscovery.Struct(&config); err != nil {
			return err
		}
		if err = setDiscoveryConfigToEnvironment(config); err != nil {
			return err
		}
		// Service.
		if err = configDiscovery.Struct(&service); err != nil {
			return err
		}
		if err = setServiceConfigToEnvironment(service); err != nil {
			return err
		}
	}
	if len(config.Endpoints) == 0 {
		g.Log().Fatal(ctx, `endpoints configuration not found for discovery`)
	}
	// Configuration: service
	configService, err := g.Cfg().Get(ctx, configNodeNameService)
	if err != nil {
		return err
	}
	if !configService.IsNil() {
		if configService.IsSlice() {
			var (
				services []*Service
			)
			if err = configService.Structs(&services); err != nil {
				return err
			}
			for _, s := range services {
				if err = setServiceConfigToEnvironment(s); err != nil {
					return err
				}
			}
		} else {
			if err = configService.Struct(&service); err != nil {
				return err
			}
			if err = setServiceConfigToEnvironment(service); err != nil {
				return err
			}
		}
	}
	return nil
}

// setDiscoveryConfigToEnvironment sets the discovery environment value with Config object.
func setDiscoveryConfigToEnvironment(config *Config) error {
	if config == nil {
		return nil
	}
	if len(config.Endpoints) > 0 {
		if err := genv.Set(EnvEndpoints, gstr.Join(config.Endpoints, ",")); err != nil {
			return err
		}
	}
	if config.KeepAlive > 0 {
		if err := genv.Set(EnvKeepAlive, config.KeepAlive.String()); err != nil {
			return err
		}
	}
	if config.PrefixRoot != "" {
		if err := genv.Set(EnvPrefixRoot, config.PrefixRoot); err != nil {
			return err
		}
	}
	return nil
}

// setServiceConfigToEnvironment sets the service environment value with Service object.
func setServiceConfigToEnvironment(service *Service) error {
	if service == nil {
		return nil
	}
	if service.AppID != "" {
		if err := genv.Set(EnvAppID, service.AppID); err != nil {
			return err
		}
	}
	if service.Address != "" {
		if err := genv.Set(EnvAddress, service.Address); err != nil {
			return err
		}
	}
	if service.Version != "" {
		if err := genv.Set(EnvVersion, service.Version); err != nil {
			return err
		}
	}
	if service.Group != "" {
		if err := genv.Set(EnvGroup, service.Group); err != nil {
			return err
		}
	}
	if service.Deployment != "" {
		if err := genv.Set(EnvDeployment, service.Deployment); err != nil {
			return err
		}
	}
	if len(service.Metadata) > 0 {
		b, _ := gjson.Marshal(service.Metadata)
		if err := genv.Set(EnvMetadata, string(b)); err != nil {
			return err
		}
	}
	return nil
}
