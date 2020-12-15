package discovery

import (
	"encoding/json"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"sync"
	"time"
)

// etcdDiscovery is the interface Registry implements using ETCD.
type etcdDiscovery struct {
	sync.RWMutex
	etcd3Client  *etcd3.Client
	keepaliveTtl time.Duration
	etcdGrantId  etcd3.LeaseID
}

var (
	// defaultDiscovery is the default Registry object that used in package method for convenience.
	defaultDiscovery Discovery
)

// SetDefault sets the default Discovery implements as your own implemented interface.
// This configuration function should be called before using function `Register`.
func SetDefault(discovery Discovery) {
	defaultDiscovery = discovery
}

// initDefaultDiscovery lazily initializes the local register object.
func initDefaultDiscovery() error {
	if defaultDiscovery != nil {
		return nil
	}
	client, err := getEtcdClient()
	if err != nil {
		return err
	}
	defaultDiscovery = &etcdDiscovery{
		etcd3Client:  client,
		keepaliveTtl: gcmd.GetWithEnv(EnvKey.KeepAlive, DefaultValue.KeepAlive).Duration(),
	}
	return nil
}

// Register registers `service` to ETCD.
func Register(service *Service) error {
	if err := initDefaultDiscovery(); err != nil {
		return err
	}
	return defaultDiscovery.Register(service)
}

// Services returns all registered service list.
func Services() []*Service {
	return defaultDiscovery.Services()
}

// Unregister removes `service` from ETCD.
func Unregister(service *Service) error {
	if err := initDefaultDiscovery(); err != nil {
		return err
	}
	return defaultDiscovery.Unregister(service)
}

// Close closes the default Registry for gracefully shutdown purpose.
func Close() error {
	if err := initDefaultDiscovery(); err != nil {
		return err
	}
	return defaultDiscovery.Close()
}

// Register registers `service` to ETCD.
func (r *etcdDiscovery) Register(service *Service) error {
	// Necessary.
	if service.AppId == "" {
		service.AppId = gcmd.GetWithEnv(EnvKey.AppId).String()
		if service.AppId == "" {
			return gerror.New(`service app id cannot be empty`)
		}
	}
	// Necessary.
	if service.Address == "" {
		service.Address = gcmd.GetWithEnv(EnvKey.Address).String()
		if service.Address == "" {
			return gerror.Newf(`service address for "%s" cannot be empty`, service.AppId)
		}
	}
	if service.Deployment == "" {
		service.Deployment = gcmd.GetWithEnv(EnvKey.Deployment, DefaultValue.Deployment).String()
	}
	if service.Group == "" {
		service.Group = gcmd.GetWithEnv(EnvKey.Group, DefaultValue.Group).String()
	}
	if service.Version == "" {
		service.Version = gcmd.GetWithEnv(EnvKey.Version, DefaultValue.Version).String()
	}
	if len(service.Metadata) == 0 {
		service.Metadata = gcmd.GetWithEnv(EnvKey.Metadata).Map()
	}
	metadataMarshalBytes, err := json.Marshal(service.Metadata)
	if err != nil {
		return err
	}
	var (
		metadataMarshalStr = string(metadataMarshalBytes)
		serviceRegisterKey = service.RegisterKey()
	)

	//g.Log().Debugf(`register key: %s`, serviceRegisterKey)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := r.etcd3Client.Grant(ctx, int64(r.keepaliveTtl/time.Second))
	if err != nil {
		return err
	}
	//g.Log().Debugf(`registered: %d, %s`, resp.ID, metadataMarshalStr)
	r.etcdGrantId = resp.ID
	if _, err := r.etcd3Client.Put(context.Background(), serviceRegisterKey, metadataMarshalStr, etcd3.WithLease(r.etcdGrantId)); err != nil {
		return err
	}
	//g.Log().Debugf(`request keepalive for grant id: %d`, resp.ID)
	keepAliceCh, err := r.etcd3Client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	g.Log().Printf(`service registered: %+v`, service)
	go r.keepAlive(service, keepAliceCh)
	return nil
}

// keepAlive continuously keeps alive the lease from ETCD.
func (r *etcdDiscovery) keepAlive(service *Service, keepAliceCh <-chan *etcd3.LeaseKeepAliveResponse) {
	for {
		select {
		case <-r.etcd3Client.Ctx().Done():
			g.Log().Debugf("keepalive done for lease id: %d", r.etcdGrantId)
			return

		case res, ok := <-keepAliceCh:
			if res != nil {
				//g.Log().Debugf(`keepalive loop: %v, %s`, ok, res.String())
			}
			if !ok {
				//g.Log().Debugf(`keepalive exit, lease id: %d`, r.etcdGrantId)
				return
			}
		}
	}
}

// Services returns all registered service list.
func (r *etcdDiscovery) Services() []*Service {
	return nil
}

// Unregister removes `service` from ETCD.
func (r *etcdDiscovery) Unregister(service *Service) error {
	//g.Log().Debugf(`discovery.Unregister: %s`, service.AppId)
	_, err := r.etcd3Client.Revoke(context.Background(), r.etcdGrantId)
	return err
}

// Close closes the Registry for gracefully shutdown purpose.
func (r *etcdDiscovery) Close() error {
	return r.etcd3Client.Close()
}
