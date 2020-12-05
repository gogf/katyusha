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

// etcdBuilder implements interface resolver.Builder.
type etcdBuilder struct {
	etcdWatcher *etcdWatcher
	waitGroup   sync.WaitGroup // Used for gracefully close the builder.
}

func init() {
	// It uses default builder handling the DNS for grpc service requests.
	resolver.Register(&etcdBuilder{})
}

// Build implements interface google.golang.org/grpc/resolver.Builder.
func (r *etcdBuilder) Build(target resolver.Target, clientConn resolver.ClientConn, options resolver.BuildOptions) (resolver.Resolver, error) {
	g.Log().Debug("Build", target, clientConn, options)
	if target.Endpoint == "" {
		return nil, gerror.New(`requested app id cannot be empty`)
	}
	endpoints := gstr.SplitAndTrim(gcmd.GetWithEnv(EnvKeyEndpoints).String(), ",")
	if len(endpoints) == 0 {
		return nil, gerror.New(`discovery server endpoints not found from environment or command-line`)
	}
	// ETCD watcher initialization.
	if r.etcdWatcher == nil {
		etcdClient, err := etcd3.New(etcd3.Config{
			Endpoints: endpoints,
		})
		if err != nil {
			return nil, err
		}
		// Just watching certain `deployment` and `group` which are the same with current client.
		// It ignores other deployment and group applications.
		r.etcdWatcher = newEtcdWatcher(
			etcdClient,
			gstr.Join([]string{
				gcmd.GetWithEnv(EnvKeyPrefixRoot, DefaultPrefixRoot).String(),
				gcmd.GetWithEnv(EnvKeyDeployment, DefaultDeployment).String(),
				gcmd.GetWithEnv(EnvKeyGroup, DefaultGroup).String(),
			}, "/"),
		)
	}
	r.waitGroup.Add(1)
	go func() {
		defer r.waitGroup.Done()
		for addresses := range r.etcdWatcher.Watch(target.Endpoint) {
			g.Log().Debugf(`AppId: %s, UpdateState: %v`, target.Endpoint, addresses)
			if len(addresses) > 0 {
				clientConn.UpdateState(resolver.State{
					Addresses: addresses,
				})
			}
		}
	}()
	return r, nil
}

// Scheme implements interface google.golang.org/grpc/resolver.Builder.
func (r *etcdBuilder) Scheme() string {
	return DefaultScheme
}

// ResolveNow implements interface google.golang.org/grpc/resolver.Resolver.
func (r *etcdBuilder) ResolveNow(o resolver.ResolveNowOptions) {

}

// Close implements interface google.golang.org/grpc/resolver.Resolver.
func (r *etcdBuilder) Close() {
	r.etcdWatcher.Close()
	r.waitGroup.Wait()
}
