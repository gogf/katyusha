// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package discovery

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/text/gstr"
)

// RegisterKey formats the service information with `registryDir` and returns the key string
// for registering.
func (s *Service) RegisterKey() string {
	return gstr.Join([]string{
		gcmd.GetOptWithEnv(EnvPrefixRoot, DefaultPrefixRoot).String(),
		s.Deployment,
		s.Group,
		s.AppID,
		s.Version,
		s.Address,
	}, "/")
}

// newServiceFromKeyValue creates and returns service from `key` and `value`.
func newServiceFromKeyValue(key, value []byte) *Service {
	array := gstr.SplitAndTrim(string(key), "/")
	if len(array) < 6 {
		return nil
	}
	service := &Service{
		Deployment: array[1],
		Group:      array[2],
		AppID:      array[3],
		Version:    array[4],
		Address:    array[5],
		Metadata:   make(g.Map),
	}
	if len(value) > 0 {
		err := gjson.DecodeTo(value, &service.Metadata)
		if err != nil {
			g.Log().Error(context.TODO(), err)
		}
	}
	return service
}
