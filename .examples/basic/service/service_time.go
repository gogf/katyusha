package service

import (
	"fmt"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/katyusha/.examples/basic/protocol"
	"golang.org/x/net/context"
)

type Time struct{}

func (s *Time) Now(ctx context.Context, r *protocol.NowReq) (*protocol.NowRes, error) {
	text := fmt.Sprintf(`%s: %s`, gcmd.GetOpt("node"), gtime.Now().String())
	return &protocol.NowRes{Time: text}, nil
}
