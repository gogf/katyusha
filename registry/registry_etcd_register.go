package registry

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"sync"
	"time"
)

type EtcdRegister struct {
	sync.RWMutex
	config      *EtcdConfig
	etcd3Client *etcd3.Client
}

type EtcdConfig struct {
	EtcdConfig  *etcd3.Config
	RegistryDir string
	TTL         time.Duration
}

func NewRegister(config *EtcdConfig) (Register, error) {
	client, err := etcd3.New(*config.EtcdConfig)
	if err != nil {
		return nil, err
	}
	registry := &EtcdRegister{
		etcd3Client: client,
		config:      config,
	}
	return registry, nil
}

func (r *EtcdRegister) Register(service *Service) error {
	serviceMarshalBytes, err := json.Marshal(service)
	if err != nil {
		return err
	}
	if service.Deployment == "" {
		service.Deployment = DeploymentDefault
	}
	if service.Group == "" {
		service.Group = DefaultGroup
	}
	var (
		ctx, _             = context.WithTimeout(context.Background(), time.Second*10)
		serviceMarshalStr  = string(serviceMarshalBytes)
		serviceRegisterKey = service.RegisterKey(r.config.RegistryDir)
	)

	g.Log().Debugf(`register key: %s`, serviceRegisterKey)
	resp, err := r.etcd3Client.Grant(ctx, int64(r.config.TTL/time.Second))
	if err != nil {
		return err
	}
	g.Log().Debugf(`registered grant id: %d`, resp.ID)
	service.etcdGrantId = resp.ID
	if _, err := r.etcd3Client.Put(ctx, serviceRegisterKey, serviceMarshalStr, etcd3.WithLease(service.etcdGrantId)); err != nil {
		return err
	}
	g.Log().Debugf(`request keepalive for grant id: %d`, resp.ID)
	keepAliceCh, err := r.etcd3Client.KeepAlive(ctx, resp.ID)
	if err != nil {
		return err
	}
	go r.keepAlive(service, keepAliceCh)
	return nil
}

func (r *EtcdRegister) keepAlive(service *Service, keepAliceCh <-chan *etcd3.LeaseKeepAliveResponse) {
	for {
		select {
		case <-r.etcd3Client.Ctx().Done():
			g.Log().Debugf("keepalive done for grant id: %d", service.etcdGrantId)
			return

		case res, ok := <-keepAliceCh:
			if res != nil {
				g.Log().Debugf(`keepalive loop: %v, %s`, ok, res.String())
			}
			if !ok {
				// I do not care about the returned error.
				r.Unregister(service)
				if err := r.Register(service); err != nil {
					g.Log().Error(err)
				}
				return
			}
		}
	}
}

func (r *EtcdRegister) Unregister(service *Service) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	_, err := r.etcd3Client.Revoke(ctx, service.etcdGrantId)
	return err
}

func (r *EtcdRegister) Close() error {
	return r.etcd3Client.Close()
}
