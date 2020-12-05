package discovery

import (
	"encoding/json"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/text/gstr"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"sync"
	"time"
)

// etcdRegister is the interface Registry implements using ETCD.
type etcdRegister struct {
	sync.RWMutex
	etcd3Client  *etcd3.Client
	keepaliveTtl time.Duration
	etcdGrantId  etcd3.LeaseID
}

var (
	// defaultRegistry is the default Registry object that used in package method for convenience.
	defaultRegistry Registry
)

// initDefaultRegister lazily initializes the local register object.
func initDefaultRegister() error {
	if defaultRegistry != nil {
		return nil
	}
	endpoints := gstr.SplitAndTrim(gcmd.GetWithEnv(EnvKeyEndpoints).String(), ",")
	if len(endpoints) == 0 {
		return gerror.New(`endpoints not found from environment or command-line`)
	}
	client, err := etcd3.New(etcd3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		return err
	}
	defaultRegistry = &etcdRegister{
		etcd3Client:  client,
		keepaliveTtl: gcmd.GetWithEnv("", DefaultKeepAliveTtl).Duration(),
	}
	return nil
}

// Register registers `service` to ETCD.
func Register(service *Service) error {
	if err := initDefaultRegister(); err != nil {
		return err
	}
	return defaultRegistry.Register(service)
}

// Unregister removes `service` from ETCD.
func Unregister(service *Service) error {
	if err := initDefaultRegister(); err != nil {
		return err
	}
	return defaultRegistry.Unregister(service)
}

// Close closes the Registry for gracefully shutdown purpose.
func Close() error {
	if err := initDefaultRegister(); err != nil {
		return err
	}
	return defaultRegistry.Close()
}

// Register registers `service` to ETCD.
func (r *etcdRegister) Register(service *Service) error {
	// Necessary.
	if service.AppId == "" {
		service.AppId = gcmd.GetWithEnv(EnvKeyAppId).String()
		if service.AppId == "" {
			return gerror.New(`service app id cannot be empty`)
		}
	}
	// Necessary.
	if service.Address == "" {
		service.Address = gcmd.GetWithEnv(EnvKeyAddress).String()
		if service.Address == "" {
			return gerror.Newf(`service address for "%s" cannot be empty`, service.AppId)
		}
	}
	if service.Deployment == "" {
		service.Deployment = gcmd.GetWithEnv(EnvKeyDeployment, DefaultDeployment).String()
	}
	if service.Group == "" {
		service.Group = gcmd.GetWithEnv(EnvKeyGroup, DefaultGroup).String()
	}
	if service.Version == "" {
		service.Version = gcmd.GetWithEnv(EnvKeyVersion, DefaultVersion).String()
	}
	metadataMarshalBytes, err := json.Marshal(service.Metadata)
	if err != nil {
		return err
	}
	var (
		metadataMarshalStr = string(metadataMarshalBytes)
		serviceRegisterKey = service.RegisterKey()
	)

	g.Log().Debugf(`register key: %s`, serviceRegisterKey)
	resp, err := r.etcd3Client.Grant(context.Background(), int64(r.keepaliveTtl/time.Second))
	if err != nil {
		return err
	}
	g.Log().Debugf(`registered: %d, %s`, resp.ID, metadataMarshalStr)
	r.etcdGrantId = resp.ID
	if _, err := r.etcd3Client.Put(context.Background(), serviceRegisterKey, metadataMarshalStr, etcd3.WithLease(r.etcdGrantId)); err != nil {
		return err
	}
	g.Log().Debugf(`request keepalive for grant id: %d`, resp.ID)
	keepAliceCh, err := r.etcd3Client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	go r.keepAlive(service, keepAliceCh)
	return nil
}

// keepAlive continuously keeps alive the lease from ETCD.
func (r *etcdRegister) keepAlive(service *Service, keepAliceCh <-chan *etcd3.LeaseKeepAliveResponse) {
	for {
		select {
		case <-r.etcd3Client.Ctx().Done():
			g.Log().Debugf("keepalive done for grant id: %d", r.etcdGrantId)
			return

		case res, ok := <-keepAliceCh:
			if res != nil {
				g.Log().Debugf(`keepalive loop: %v, %s`, ok, res.String())
			}
			if !ok {
				g.Log().Debugf(`keepalive Unregister: %s`, r.etcdGrantId)
				r.Unregister(service)
				if err := r.Register(service); err != nil {
					g.Log().Error(err)
				}
				return
			}
		}
	}
}

// Unregister removes `service` from ETCD.
func (r *etcdRegister) Unregister(service *Service) error {
	_, err := r.etcd3Client.Revoke(context.Background(), r.etcdGrantId)
	return err
}

// Close closes the Registry for gracefully shutdown purpose.
func (r *etcdRegister) Close() error {
	return r.etcd3Client.Close()
}
