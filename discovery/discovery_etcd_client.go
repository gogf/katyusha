// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package discovery

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/text/gstr"
	etcd3 "go.etcd.io/etcd/client/v3"
)

var (
	// etcdClient is the client instance for etcd.
	etcdClient *etcd3.Client
)

// getEtcdClient creates and returns an instance for etcd client.
// It returns the same instance object if it already created one.
func getEtcdClient() (*etcd3.Client, error) {
	if etcdClient != nil {
		return etcdClient, nil
	}
	endpoints := gstr.SplitAndTrim(gcmd.GetOptWithEnv(EnvEndpoints).String(), ",")
	if len(endpoints) == 0 {
		return nil, gerror.New(`endpoints not found from environment, command-line or configuration file`)
	}
	client, err := etcd3.New(etcd3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, err
	}
	etcdClient = client
	return etcdClient, nil
}
