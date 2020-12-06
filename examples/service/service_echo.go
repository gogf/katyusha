package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/katyusha/examples/protocol"
	"golang.org/x/net/context"
)

type Echo struct{}

func (s *Echo) Say(ctx context.Context, r *protocol.SayReq) (*protocol.SayRes, error) {
	g.Log().Println("Received:", r.Content)
	text := fmt.Sprintf(`%s: > %s`, gcmd.GetOpt("node"), r.Content)
	return &protocol.SayRes{Content: text}, nil
}
