package registry

import (
	etcd3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
)

// EtcdResolver implements interface resolver.Builder.
type EtcdResolver struct {
	EtcdScheme  string      // Scheme returns the scheme supported by this resolver.
	EtcdConfig  *EtcdConfig // ETCD configuration object.
	Service     *Service    // Service configuration.
	etcdWatcher *EtcdWatcher
	clientConn  resolver.ClientConn
	waitGroup   sync.WaitGroup
}

// Build implements interface google.golang.org/grpc/resolver.Builder.
func (r *EtcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	etcdClient, err := etcd3.New(*r.EtcdConfig.EtcdConfig)
	if err != nil {
		return nil, err
	}
	r.clientConn = cc
	r.etcdWatcher = newEtcdWatcher(r.Service.RegisterKey(r.EtcdConfig.RegistryDir), etcdClient)
	r.waitGroup.Add(1)
	go func() {
		defer r.waitGroup.Done()
		for addr := range r.etcdWatcher.Watch() {
			r.clientConn.UpdateState(resolver.State{Addresses: addr})
		}
	}()
	return r, nil
}

// Scheme implements interface google.golang.org/grpc/resolver.Builder.
func (r *EtcdResolver) Scheme() string {
	return r.EtcdScheme
}

// ResolveNow implements interface google.golang.org/grpc/resolver.Resolver.
func (r *EtcdResolver) ResolveNow(o resolver.ResolveNowOptions) {

}

// Close implements interface google.golang.org/grpc/resolver.Resolver.
func (r *EtcdResolver) Close() {
	r.etcdWatcher.Close()
	r.waitGroup.Wait()
}
