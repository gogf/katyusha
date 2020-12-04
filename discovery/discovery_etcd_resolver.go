package discovery

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/text/gstr"
	etcd3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
)

// EtcdResolver implements interface resolver.Builder.
type EtcdResolver struct {
	etcdWatcher *EtcdWatcher
	clientConn  resolver.ClientConn
	waitGroup   sync.WaitGroup
}

// Build implements interface google.golang.org/grpc/resolver.Builder.
func (r *EtcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	endpoints := gstr.SplitAndTrim(gcmd.GetWithEnv(EnvKeyEndpoints).String(), ",")
	if len(endpoints) == 0 {
		return nil, gerror.New(`endpoints not found from environment or command-line`)
	}
	etcdClient, err := etcd3.New(etcd3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, err
	}
	r.clientConn = cc
	r.etcdWatcher = newEtcdWatcher(
		gcmd.GetWithEnv(EnvKeyPrefixRoot, DefaultPrefixRoot).String(),
		etcdClient,
	)
	r.waitGroup.Add(1)
	go func() {
		defer r.waitGroup.Done()
		for address := range r.etcdWatcher.Watch() {
			g.Log().Debugf(`UpdateState: %v`, address)
			r.clientConn.UpdateState(resolver.State{Addresses: address})
		}
	}()
	return r, nil
}

// Scheme implements interface google.golang.org/grpc/resolver.Builder.
func (r *EtcdResolver) Scheme() string {
	return DefaultScheme
}

// ResolveNow implements interface google.golang.org/grpc/resolver.Resolver.
func (r *EtcdResolver) ResolveNow(o resolver.ResolveNowOptions) {

}

// Close implements interface google.golang.org/grpc/resolver.Resolver.
func (r *EtcdResolver) Close() {
	r.etcdWatcher.Close()
	r.waitGroup.Wait()
}
