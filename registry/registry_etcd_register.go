package registry

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
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
	EtcdConfig  etcd3.Config
	RegistryDir string
	TTL         time.Duration
}

func NewRegister(config *EtcdConfig) (Register, error) {
	client, err := etcd3.New(config.EtcdConfig)
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
	var (
		ctx                = context.Background()
		serviceMarshalStr  = string(serviceMarshalBytes)
		serviceRegisterKey = gstr.Join([]string{
			r.config.RegistryDir,
			service.Environment,
			service.Group,
			service.Version,
			service.AppId,
		}, "/")
	)

	resp, err := r.etcd3Client.Grant(ctx, int64(r.config.TTL/time.Second))
	if err != nil {
		return err
	}
	service.etcdGrantId = resp.ID
	if _, err := r.etcd3Client.Put(ctx, serviceRegisterKey, serviceMarshalStr, etcd3.WithLease(service.etcdGrantId)); err != nil {
		return err
	}

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
			return

		case _, ok := <-keepAliceCh:
			if !ok {
				if err := r.Unregister(service); err != nil {
					g.Log().Error(err)
				}
				if err := r.Register(service); err != nil {
					g.Log().Error(err)
				}
				return
			}
		}
	}
}

func (r *EtcdRegister) Unregister(service *Service) error {
	_, err := r.etcd3Client.Revoke(context.Background(), service.etcdGrantId)
	return err
}

func (r *EtcdRegister) Close() error {
	return r.etcd3Client.Close()
}
