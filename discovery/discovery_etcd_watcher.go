package discovery

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"sync"
)

type EtcdWatcher struct {
	key        string             // Watched key prefix.
	ctx        context.Context    // Context for request handling.
	cancelFunc context.CancelFunc // Cancel function for this context.
	etcdClient *etcd3.Client      // ETCD client.
	waitGroup  sync.WaitGroup     // WaitGroup for gracefully closing..
	addressMap *gmap.StrAnyMap    // Service AppId to its address list mapping, type: map[string][]resolver.Address.
}

func newEtcdWatcher(key string, etcdClient *etcd3.Client) *EtcdWatcher {
	ctx, cancelFunc := context.WithCancel(context.Background())
	w := &EtcdWatcher{
		key:        key,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		etcdClient: etcdClient,
		addressMap: gmap.NewStrAnyMap(true),
	}
	return w
}

// Watch keeps watching the registered prefix key events.
func (w *EtcdWatcher) Watch(appId string) chan []resolver.Address {
	w.initialize()
	addressCh := make(chan []resolver.Address, 10)
	w.waitGroup.Add(1)
	go func() {
		defer func() {
			close(addressCh)
			w.waitGroup.Done()
		}()
		w.updateAddressesToCh(appId, addressCh)
		// Watch events handling.
		for watchResponse := range w.etcdClient.Watch(w.ctx, w.key, etcd3.WithPrefix()) {
			for _, ev := range watchResponse.Events {
				g.Log().Debugf("watch event: %d, key: %s, value: %s", ev.Type, ev.Kv.Key, ev.Kv.Value)
				switch ev.Type {
				case mvccpb.PUT:
					service := newServiceFromKeyValue(ev.Kv.Key, ev.Kv.Value)
					if service == nil {
						g.Log().Error("service creating failed for key: %s, value:%s", ev.Kv.Key, ev.Kv.Value)
						continue
					}
					address := resolver.Address{
						Addr:       service.Address,
						Attributes: attributes.New(gconv.Interfaces(service.Metadata)...),
					}
					if w.addAddress(service.AppId, address) {
						if service.AppId == appId {
							w.updateAddressesToCh(service.AppId, addressCh)
						}
					}

				case mvccpb.DELETE:
					service := newServiceFromKeyValue(ev.Kv.Key, ev.Kv.Value)
					if service == nil {
						g.Log().Error("service creating failed for key: %s, value:%s", ev.Kv.Key, ev.Kv.Value)
						continue
					}
					address := resolver.Address{
						Addr:       service.Address,
						Attributes: attributes.New(gconv.Interfaces(service.Metadata)...),
					}
					if w.removeAddress(service.AppId, address) {
						if service.AppId == appId {
							w.updateAddressesToCh(service.AppId, addressCh)
						}
					}
				}
			}
		}
	}()
	return addressCh
}

// Initialize retrieves data from discovery server and initializes the address map.
func (w *EtcdWatcher) initialize() {
	w.addressMap.LockFunc(func(m map[string]interface{}) {
		if len(m) == 0 {
			res, err := w.etcdClient.Get(w.ctx, w.key, etcd3.WithPrefix())
			if err == nil {
				services := extractServices(res)
				if len(services) > 0 {
					for _, service := range services {
						if _, ok := m[service.AppId]; !ok {
							m[service.AppId] = make([]resolver.Address, 0)
						}
						m[service.AppId] = append(m[service.AppId].([]resolver.Address), resolver.Address{
							Addr:       service.Address,
							Attributes: attributes.New(gconv.Interfaces(service.Metadata)...),
						})
					}
				}
			}
		}
	})
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
	g.Log().Debugf(`extractServices: %v`, services)
	return services
}

func (w *EtcdWatcher) updateAddressesToCh(appId string, addressCh chan []resolver.Address) {
	clonedAddresses := make([]resolver.Address, 0)
	w.addressMap.RLockFunc(func(m map[string]interface{}) {
		if v := m[appId]; v == nil {
			return
		} else {
			for _, address := range v.([]resolver.Address) {
				clonedAddresses = append(clonedAddresses, address)
			}
		}
	})
	if len(clonedAddresses) > 0 {
		addressCh <- clonedAddresses
	}
}

func (w *EtcdWatcher) addAddress(appId string, address resolver.Address) bool {
	w.addressMap.LockFunc(func(m map[string]interface{}) {
		if _, ok := m[appId]; !ok {
			m[appId] = make([]resolver.Address, 0)
		}
		for _, v := range m[appId].([]resolver.Address) {
			if address.Addr == v.Addr {
				// Already added.
				return
			}
		}
		m[appId] = append(m[appId].([]resolver.Address), address)
	})
	return true
}

func (w *EtcdWatcher) removeAddress(appId string, address resolver.Address) bool {
	w.addressMap.LockFunc(func(m map[string]interface{}) {
		if _, ok := m[appId]; !ok {
			m[appId] = make([]resolver.Address, 0)
		}
		for i, v := range m[appId].([]resolver.Address) {
			if address.Addr == v.Addr {
				m[appId] = append(
					m[appId].([]resolver.Address)[:i],
					m[appId].([]resolver.Address)[i+1:]...,
				)
				return
			}
		}
	})
	return false
}

func (w *EtcdWatcher) Close() {
	w.cancelFunc()
}
