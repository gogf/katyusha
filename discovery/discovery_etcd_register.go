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

type EtcdRegister struct {
	sync.RWMutex
	etcd3Client  *etcd3.Client
	keepaliveTtl time.Duration
	etcdGrantId  etcd3.LeaseID
}

func NewRegister() (Register, error) {
	endpoints := gstr.SplitAndTrim(gcmd.GetWithEnv(EnvKeyEndpoints).String(), ",")
	if len(endpoints) == 0 {
		return nil, gerror.New(`endpoints not found from environment or command-line`)
	}
	client, err := etcd3.New(etcd3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, err
	}
	registry := &EtcdRegister{
		etcd3Client:  client,
		keepaliveTtl: gcmd.GetWithEnv("", DefaultKeepAliveTtl).Duration(),
	}
	return registry, nil
}

func (r *EtcdRegister) Register(service *Service) error {
	if service.Deployment == "" {
		service.Deployment = gcmd.GetWithEnv(EnvKeyDeployment, DefaultDeployment).String()
	}
	if service.Group == "" {
		service.Group = gcmd.GetWithEnv(EnvKeyGroup, DefaultGroup).String()
	}
	if service.Version == "" {
		service.Version = gcmd.GetWithEnv(EnvKeyVersion, DefaultVersion).String()
	}
	serviceMarshalBytes, err := json.Marshal(service)
	if err != nil {
		return err
	}
	var (
		serviceMarshalStr  = string(serviceMarshalBytes)
		serviceRegisterKey = service.RegisterKey()
	)

	g.Log().Debugf(`register key: %s`, serviceRegisterKey)
	resp, err := r.etcd3Client.Grant(context.Background(), int64(r.keepaliveTtl/time.Second))
	if err != nil {
		return err
	}
	g.Log().Debugf(`registered: %d, %s`, resp.ID, serviceMarshalStr)
	r.etcdGrantId = resp.ID
	if _, err := r.etcd3Client.Put(context.Background(), serviceRegisterKey, serviceMarshalStr, etcd3.WithLease(r.etcdGrantId)); err != nil {
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

func (r *EtcdRegister) keepAlive(service *Service, keepAliceCh <-chan *etcd3.LeaseKeepAliveResponse) {
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

func (r *EtcdRegister) Unregister(service *Service) error {
	_, err := r.etcd3Client.Revoke(context.Background(), r.etcdGrantId)
	return err
}

func (r *EtcdRegister) Close() error {
	return r.etcd3Client.Close()
}
