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
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 退出信号
	port   string             // 服务的端口
}

// 启动服务
func NewVortexCore(ctx context.Context, opts ...Option) error {
	opt := &Vortex{}
	for _, o := range opts {
		o(opt)
	}
	return nil
}

// 开启端口监听，先判断当前请求的协议，然后选择对应的协议进行处理
func Start() {
}
