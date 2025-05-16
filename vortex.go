package Vortex

import "context"

type Option func(*Vortex)

func WithListenPort(port string) Option {
	return func(v *Vortex) {
		v.port = port
	}
}

// 框架的整体结构
type Vortex struct {
	ctx  context.Context
	port string // 服务的端口
}

// 启动服务
func Start(ctx context.Context, opts ...Option) error {
	opt := &Vortex{}
	for _, o := range opts {
		o(opt)
	}
	return nil
}
