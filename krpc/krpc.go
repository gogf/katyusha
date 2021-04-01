// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package krpc

import (
	"github.com/gogf/katyusha/krpc/internal/grpcctx"
)

type (
	krpcClient struct{}
	krpcServer struct{}
)

const (
	configNodeNameGrpcServer = "grpcserver"
	configNodeNameHttpServer = "httpserver"
	randomPortMin            = 8000
	randomPortMax            = 9999
	randomPortNotAvailable   = -1
)

var (
	Ctx    = grpcctx.Ctx  // Ctx manages the context feature.
	Client = krpcClient{} // Client manages the client features.
	Server = krpcServer{} // Server manages the server feature.
)
