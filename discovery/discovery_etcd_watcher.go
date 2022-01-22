// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package discovery

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gutil"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

type etcdWatcher struct {
	key        string             // Watched key prefix.
	ctx        context.Context    // Context for request handling.
	cancelFunc context.CancelFunc // Cancel function for this context.
	etcdClient *etcd3.Client      // ETCD client.
	waitGroup  sync.WaitGroup     // WaitGroup for gracefully closing..
	addresses  []resolver.Address // Address list for certain service.
}

func newEtcdWatcher(etcdClient *etcd3.Client, key string) *etcdWatcher {
	ctx, cancelFunc := context.WithCancel(context.Background())
	w := &etcdWatcher{
		key:        key,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		etcdClient: etcdClient,
		addresses:  make([]resolver.Address, 0),
	}
	return w
}

// Watch keeps watching the registered prefix key events.
func (w *etcdWatcher) Watch() chan []resolver.Address {
	var ctx = context.TODO()
	w.initializeAddresses()
	addressCh := make(chan []resolver.Address, 10)
	w.waitGroup.Add(1)
	go func() {
		defer func() {
			close(addressCh)
			w.waitGroup.Done()
		}()
		w.updateAddressesToCh(addressCh)
		// Watch events handling.
		for watchResponse := range w.etcdClient.Watch(w.ctx, w.key, etcd3.WithPrefix()) {
			for _, ev := range watchResponse.Events {
				g.Log().Debugf(ctx, "watch event: %d, key: %s, value: %s", ev.Type, ev.Kv.Key, ev.Kv.Value)
				switch ev.Type {
				case mvccpb.PUT:
					service := newServiceFromKeyValue(ev.Kv.Key, ev.Kv.Value)
					if service == nil {
						g.Log().Error(ctx, "service creating failed for key: %s, value:%s", ev.Kv.Key, ev.Kv.Value)
						continue
					}
					address := resolver.Address{
						Addr:       service.Address,
						Attributes: attributes.New(gutil.MapToSlice(service.Metadata)...),
					}
					if w.addAddress(address) {
						w.updateAddressesToCh(addressCh)
					}

				case mvccpb.DELETE:
					service := newServiceFromKeyValue(ev.Kv.Key, ev.Kv.Value)
					if service == nil {
						g.Log().Error(ctx, "service creating failed for key: %s, value:%s", ev.Kv.Key, ev.Kv.Value)
						continue
					}
					address := resolver.Address{
						Addr:       service.Address,
						Attributes: attributes.New(gutil.MapToSlice(service.Metadata)...),
					}
					if w.removeAddress(address) {
						w.updateAddressesToCh(addressCh)
					}
				}
			}
		}
	}()
	return addressCh
}

// initializeAddresses retrieves data from discovery server and initializes the address list.
func (w *etcdWatcher) initializeAddresses() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	res, err := w.etcdClient.Get(ctx, w.key, etcd3.WithPrefix())
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if res == nil {
		return
	}
	services := extractServices(res)
	if len(services) > 0 {
		w.addresses = make([]resolver.Address, 0)
		for _, service := range services {
			w.addresses = append(w.addresses, resolver.Address{
				Addr:       service.Address,
				Attributes: attributes.New(gutil.MapToSlice(service.Metadata)...),
			})
		}
	}
}

// extractServices extracts etcd watch response context to service list.
func extractServices(resp *etcd3.GetResponse) []*Service {
	var services []*Service
	if resp == nil || resp.Kvs == nil {
		return services
	}
	for _, kv := range resp.Kvs {
		if service := newServiceFromKeyValue(kv.Key, kv.Value); service != nil {
			services = append(services, service)
		}
	}
	// g.Log().Debugf(`extractServices: %v`, services)
	return services
}

func (w *etcdWatcher) updateAddressesToCh(addressCh chan []resolver.Address) {
	clonedAddresses := make([]resolver.Address, 0)
	for _, address := range w.addresses {
		clonedAddresses = append(clonedAddresses, address)
	}
	addressCh <- clonedAddresses
}

func (w *etcdWatcher) addAddress(address resolver.Address) bool {
	for _, v := range w.addresses {
		if address.Addr == v.Addr {
			// Already added.
			return false
		}
	}
	w.addresses = append(w.addresses, address)
	return true
}

func (w *etcdWatcher) removeAddress(address resolver.Address) bool {
	for i, v := range w.addresses {
		if address.Addr == v.Addr {
			w.addresses = append(w.addresses[:i], w.addresses[i+1:]...)
			return true
		}
	}
	return false
}

func (w *etcdWatcher) Close() {
	w.cancelFunc()
}
