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
	cc            resolver.ClientConn
	wg            sync.WaitGroup
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	etcdCli, err := etcd3.New(r.etcdConfig)
	if err != nil {
		return nil, err
	}
	r.cc = cc
	r.watcher = newWatcher(r.etcdWatchPath, etcdCli)
	r.start()
	return r, nil
}

func (r *etcdResolver) Scheme() string {
	return r.scheme
}

func (r *etcdResolver) start() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		out := r.watcher.Watch()
		for addr := range out {
			r.cc.UpdateState(resolver.State{Addresses: addr})
		}
	}()
}

func (r *etcdResolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (r *etcdResolver) Close() {
	r.watcher.Close()
	r.wg.Wait()
}

func RegisterResolver(scheme string, etcdConfig etcd3.Config, registryDir, srvName, srvVersion string) {
	resolver.Register(&etcdResolver{
		scheme:        scheme,
		etcdConfig:    etcdConfig,
		etcdWatchPath: registryDir + "/" + srvName + "/" + srvVersion,
	})
}
