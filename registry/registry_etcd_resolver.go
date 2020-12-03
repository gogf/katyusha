package registry

import (
	etcd3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
)

type etcdResolver struct {
	scheme        string
	etcdConfig    etcd3.Config
	etcdWatchPath string
	watcher       *Watcher
	clientConn    resolver.ClientConn
	waitGroup     sync.WaitGroup
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	etcdClient, err := etcd3.New(r.etcdConfig)
	if err != nil {
		return nil, err
	}
	r.clientConn = cc
	r.watcher = newWatcher(r.etcdWatchPath, etcdClient)
	r.waitGroup.Add(1)
	go func() {
		defer r.waitGroup.Done()
		for addr := range r.watcher.Watch() {
			r.clientConn.UpdateState(resolver.State{Addresses: addr})
		}
	}()
	return r, nil
}

func (r *etcdResolver) Scheme() string {
	return r.scheme
}

func (r *etcdResolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (r *etcdResolver) Close() {
	r.watcher.Close()
	r.waitGroup.Wait()
}

func RegisterResolver(scheme string, etcdConfig etcd3.Config, registryDir, srvName, srvVersion string) {
	resolver.Register(&etcdResolver{
		scheme:     scheme,
		etcdConfig: etcdConfig,
		//etcdWatchPath: registryDir + "/" + srvName + "/" + srvVersion,
		etcdWatchPath: registryDir + "/default/default/test/v1.0",
	})
}
