package registry

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
	"sync"
	"time"
)

type RegisterEtcd struct {
	sync.RWMutex
	config      *Config
	etcd3Client *etcd3.Client
	canceler    map[string]context.CancelFunc
}

type Config struct {
	EtcdConfig  etcd3.Config
	RegistryDir string
	Ttl         time.Duration
}

func NewRegister(config *Config) (Register, error) {
	client, err := etcd3.New(config.EtcdConfig)
	if err != nil {
		return nil, err
	}

	registry := &RegisterEtcd{
		etcd3Client: client,
		config:      config,
		canceler:    make(map[string]context.CancelFunc),
	}
	return registry, nil
}

func (r *RegisterEtcd) Register(service *Service) error {
	val, err := json.Marshal(service)
	if err != nil {
		return err
	}

	key := r.config.RegistryDir + "/" + service.Name + "/" + service.Version + "/" + service.AppId
	value := string(val)
	ctx, cancel := context.WithCancel(context.Background())
	r.Lock()
	r.canceler[service.AppId] = cancel
	r.Unlock()

	insertFunc := func() error {
		resp, err := r.etcd3Client.Grant(ctx, int64(r.config.Ttl/time.Second))
		if err != nil {
			fmt.Printf("[Register] %v\n", err.Error())
			return err
		}
		_, err = r.etcd3Client.Get(ctx, key)
		if err != nil {
			if err == rpctypes.ErrKeyNotFound {
				if _, err := r.etcd3Client.Put(ctx, key, value, etcd3.WithLease(resp.ID)); err != nil {
					grpclog.Infof("grpclb: set key '%s' with ttl to etcd3 failed: %s", key, err.Error())
				}
			} else {
				grpclog.Infof("grpclb: key '%s' connect to etcd3 failed: %s", key, err.Error())
			}
			return err
		} else {
			// refresh set to true for not notifying the watcher
			if _, err := r.etcd3Client.Put(ctx, key, value, etcd3.WithLease(resp.ID)); err != nil {
				grpclog.Infof("grpclb: refresh key '%s' with ttl to etcd3 failed: %s", key, err.Error())
				return err
			}
		}
		return nil
	}

	err = insertFunc()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(r.config.Ttl / 5)
	for {
		select {
		case <-ticker.C:
			insertFunc()
		case <-ctx.Done():
			ticker.Stop()
			if _, err := r.etcd3Client.Delete(context.Background(), key); err != nil {
				grpclog.Infof("grpclb: deregister '%s' failed: %s", key, err.Error())
			}
			return nil
		}
	}

	return nil
}

func (r *RegisterEtcd) Unregister(service *Service) error {
	r.RLock()
	cancel, ok := r.canceler[service.AppId]
	r.RUnlock()

	if ok {
		cancel()
	}
	return nil
}

func (r *RegisterEtcd) Close() {
	r.etcd3Client.Close()
}
