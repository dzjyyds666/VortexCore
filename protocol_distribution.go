package Vortex

import (
	"bytes"
	"context"
	"net"
)

// 根据请求的前几个字节做协议识别，并并发连接
// 分发器
type Dispatcher struct {
	ctx context.Context
	net.Conn
	buf bytes.Buffer
}

func NewDispatcher(ctx context.Context, conn net.Conn) *Dispatcher {
	return &Dispatcher{
		ctx:  ctx,
		Conn: conn,
	}
}
