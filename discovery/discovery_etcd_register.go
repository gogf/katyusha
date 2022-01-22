// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package discovery

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)

// etcdDiscovery is the interface Registry implements using ETCD.
type etcdDiscovery struct {
	sync.RWMutex
	etcd3Client  *etcd3.Client
	keepaliveTTL time.Duration
	etcdGrantID  etcd3.LeaseID
}

var (
	// defaultDiscovery is the default Registry object that used in package method for convenience.
	defaultDiscovery Discovery
)

// Register registers `service` to ETCD.
func (r *etcdDiscovery) Register(service *Service) error {
	var ctx = context.TODO()
	// Necessary.
	if service.AppID == "" {
		service.AppID = gcmd.GetOptWithEnv(EnvAppID).String()
		if service.AppID == "" {
			return gerror.New(`service app id cannot be empty`)
		}
	}
	// Necessary.
	if service.Address == "" {
		service.Address = gcmd.GetOptWithEnv(EnvAddress).String()
		if service.Address == "" {
			return gerror.Newf(`service address for "%s" cannot be empty`, service.AppID)
		}
	}
	if service.Deployment == "" {
		service.Deployment = gcmd.GetOptWithEnv(EnvDeployment, DefaultDeployment).String()
	}
	if service.Group == "" {
		service.Group = gcmd.GetOptWithEnv(EnvGroup, DefaultGroup).String()
	}
	if service.Version == "" {
		service.Version = gcmd.GetOptWithEnv(EnvVersion, DefaultVersion).String()
	}
	if len(service.Metadata) == 0 {
		service.Metadata = gcmd.GetOptWithEnv(EnvMetadata).Map()
	}
	metadataMarshalBytes, err := gjson.Marshal(service.Metadata)
	if err != nil {
		return err
	}
	var (
		metadataMarshalStr = string(metadataMarshalBytes)
		serviceRegisterKey = service.RegisterKey()
	)

	g.Log().Debugf(ctx, `service register key: %s`, serviceRegisterKey)
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	resp, err := r.etcd3Client.Grant(ctx, int64(r.keepaliveTTL/time.Second))
	if err != nil {
		return err
	}
	g.Log().Debugf(ctx, `service registered lease id: %d, metadata: %s`, resp.ID, metadataMarshalStr)
	r.etcdGrantID = resp.ID
	if _, err = r.etcd3Client.Put(
		ctx, serviceRegisterKey, metadataMarshalStr, etcd3.WithLease(r.etcdGrantID),
	); err != nil {
		return err
	}
	g.Log().Debugf(ctx, `service request keepalive for grant id: %d`, resp.ID)
	keepAliceCh, err := r.etcd3Client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	g.Log().Printf(ctx, `service registered: %+v`, service)
	go r.keepAlive(service, keepAliceCh)
	return nil
}

// keepAlive continuously keeps alive the lease from ETCD.
func (r *etcdDiscovery) keepAlive(service *Service, keepAliceCh <-chan *etcd3.LeaseKeepAliveResponse) {
	var ctx = context.TODO()
	for {
		select {
		case <-r.etcd3Client.Ctx().Done():
			g.Log().Debugf(ctx, "keepalive done for lease id: %d", r.etcdGrantID)
			return

		case res, ok := <-keepAliceCh:
			if res != nil {
				// g.Log().Debugf(ctx, `keepalive loop: %v, %s`, ok, res.String())
			}
			if !ok {
				// g.Log().Debugf(ctx, `keepalive exit, lease id: %d`, r.etcdGrantID)
				return
			}
		}
	}
}

// Services returns all registered service list.
// TODO implements.
func (r *etcdDiscovery) Services() ([]*Service, error) {
	return nil, gerror.New("not implemented")
}

// Unregister removes `service` from ETCD.
func (r *etcdDiscovery) Unregister(service *Service) error {
	_, err := r.etcd3Client.Revoke(context.Background(), r.etcdGrantID)
	return err
}

// Close closes the Registry for gracefully shutdown purpose.
func (r *etcdDiscovery) Close() error {
	return r.etcd3Client.Close()
}
