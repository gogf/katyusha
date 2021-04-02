// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package discovery

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/text/gstr"
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
	// Create a new builder for each client.
	builder := &etcdBuilder{}
	// Etcd client instance.
	etcdClient, err := getEtcdClient()
	if err != nil {
		return nil, err
	}
	// Watch certain service prefix.
	builder.etcdWatcher = newEtcdWatcher(
		etcdClient,
		gstr.Join([]string{
			gcmd.GetOptWithEnv(EnvKey.PrefixRoot, DefaultValue.PrefixRoot).String(),
			gcmd.GetOptWithEnv(EnvKey.Deployment, DefaultValue.Deployment).String(),
			gcmd.GetOptWithEnv(EnvKey.Group, DefaultValue.Group).String(),
			target.Endpoint,
		}, "/"),
	)
	builder.waitGroup.Add(1)
	go func() {
		defer builder.waitGroup.Done()
		for addresses := range builder.etcdWatcher.Watch() {
			g.Log().Debugf(`AppId: %s, UpdateState: %v`, target.Endpoint, addresses)
			if len(addresses) > 0 {
				clientConn.UpdateState(resolver.State{
					Addresses: addresses,
				})
			} else {
				// Service addresses empty, that means service shuts down or unavailable temporarily.
				clientConn.ReportError(gerror.New("Service unavailable: service shuts down or unavailable temporarily"))
			}
		}
	}()
	return r, nil
}

// Scheme implements interface google.golang.org/grpc/resolver.Builder.
func (r *etcdBuilder) Scheme() string {
	return DefaultValue.Scheme
}

// ResolveNow implements interface google.golang.org/grpc/resolver.Resolver.
func (r *etcdBuilder) ResolveNow(opts resolver.ResolveNowOptions) {
	//g.Log().Debug("ResolveNow:", opts)
}

// Close implements interface google.golang.org/grpc/resolver.Resolver.
func (r *etcdBuilder) Close() {
	r.etcdWatcher.Close()
	r.waitGroup.Wait()
}
