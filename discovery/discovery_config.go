package discovery

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
)

func init() {
	if !g.Cfg().Available() {
		return
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
			g.Log().Error(err)
		}
		discoveryConfigToEnvironment(config)

		// Service.
		if err := configDiscovery.Struct(&service); err != nil {
			g.Log().Error(err)
		}
		serviceConfigToEnvironment(service)
	}
	// Configuration: service
	configService := g.Cfg().GetVar(configNodeNameService)
	if !configService.IsNil() {
		if configService.IsSlice() {
			var (
				services []*Service
			)
			if err := configService.Structs(&services); err != nil {
				g.Log().Error(err)
			}
			for _, service := range services {
				serviceConfigToEnvironment(service)
			}
		} else {
			var (
				service *Service
			)
			if err := configService.Struct(&service); err != nil {
				g.Log().Error(err)
			}
			serviceConfigToEnvironment(service)
		}
	}
}

// SetConfig sets the discovery configuration using Config.
func SetConfig(config *Config) {
	discoveryConfigToEnvironment(config)
}

// discoveryConfigToEnvironment sets the discovery environment value with Config object.
func discoveryConfigToEnvironment(config *Config) {
	if config == nil {
		return
	}
	if len(config.Endpoints) > 0 {
		genv.Set(EnvKey.Endpoints, gstr.Join(config.Endpoints, ","))
	}
	if config.KeepAlive > 0 {
		genv.Set(EnvKey.KeepAlive, config.KeepAlive.String())
	}
	if config.PrefixRoot != "" {
		genv.Set(EnvKey.PrefixRoot, config.PrefixRoot)
	}
}

// serviceConfigToEnvironment sets the service environment value with Service object.
func serviceConfigToEnvironment(service *Service) {
	if service == nil {
		return
	}
	if service.AppId != "" {
		genv.Set(EnvKey.AppId, service.AppId)
	}
	if service.Address != "" {
		genv.Set(EnvKey.Address, service.Address)
	}
	if service.Version != "" {
		genv.Set(EnvKey.Version, service.Version)
	}
	if service.Group != "" {
		genv.Set(EnvKey.Group, service.Group)
	}
	if service.Deployment != "" {
		genv.Set(EnvKey.Deployment, service.Deployment)
	}
	if len(service.Metadata) > 0 {
		b, _ := json.Marshal(service.Metadata)
		genv.Set(EnvKey.Metadata, string(b))
	}
}
