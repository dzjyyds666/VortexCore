package vortex

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
)

var Transport = struct {
	TCP string
	UDP string
}{
	TCP: "tcp",
	UDP: "udp",
}

type Option func(*Vortex)

func WithListenPort(port string) Option {
	return func(v *Vortex) {
		v.port = port
	}
}

func WithTransport(transport string) Option {
	return func(v *Vortex) {
		v.transport = transport
	}
}

func WithHttp1() Option {
	return func(v *Vortex) {
		v.protocol = append(v.protocol, http1)
	}
}

func WithWebSocket() Option {
	return func(v *Vortex) {
		v.protocol = append(v.protocol, webSocket)
	}
}

func WithHttp2() Option {
	return func(v *Vortex) {
		v.protocol = append(v.protocol, http2)
	}
}

func WithHttpRouter(routers []*HttpRouter) Option {
	return func(v *Vortex) {
		if v.httpRouter == nil {
			v.httpRouter = make([]*HttpRouter, 0)
		}
		v.httpRouter = append(v.httpRouter, routers...)
	}
}

// 框架的整体结构
type Vortex struct {
	ctx        context.Context    // 上下文
	cancel     context.CancelFunc // 退出信号
	port       string             // 服务的端口
	transport  string             // 传输协议
	protocol   []string           // 支持的协议列表
	httpServ   *httpServer        // http服务，封装了echo框架
	httpRouter []*HttpRouter      // http服务路由表
}

// 启动服务
func NewVortexCore(ctx context.Context, opts ...Option) *Vortex {
	vortex := &Vortex{
		transport: Transport.TCP, // 默认使用 TCP
	}
	for _, o := range opts {
		o(vortex)
	}

	if len(vortex.port) <= 0 {
		panic("port must be set")
	}

	for _, p := range vortex.protocol {
		switch p {
		case http1:
			// 添加默认的Http路由
			vortex.httpRouter = append(vortex.httpRouter,
				AppendHttpRouter([]string{http.MethodGet}, "/checkAlive", func(ctx VortexContext) error {
					return ctx.GetEcho().String(http.StatusOK, "ok")
				}))
			NewHttpServer(ctx, vortex.httpRouter)
		}
	}

	return vortex
}

// 开启端口监听，先判断当前请求的协议，然后选择对应的协议进行处理
func (v *Vortex) BootStorp() {
	ln, err := net.Listen(v.transport, fmt.Sprintf(":%s", v.port))
	if nil != err {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if nil != err {
			fmt.Printf("accept error: %v\n", err)
			continue
		}
		go v.ParsingRequest(conn) // 异步处理请求
	}
}

func (v *Vortex) ParsingRequest(conn net.Conn) {
	// 这里可以实现协议解析逻辑
	// 例如读取前几个字节来判断是 HTTP 还是 WebSocket 等
	// 然后根据协议类型进行相应的处理
	ctx, cancel := context.WithCancel(v.ctx)
	defer func() {
		cancel()
		conn.Close()
	}()
	d := NewDispatcher(ctx, conn)
	protocl, err := d.Parse()
	// 关闭连接
	if nil != err || protocl == "unknown" {
		fmt.Printf("parse error: %v\n", err)
		return
	}

	switch protocl {
	case http1:
		// 使用echo框架处理 HTTP/1.1 请求
		v.handleHttpWithEcho(d)
	case webSocket:
		// 使用 WebSocket 处理逻辑
	case http2:
		// 使用 HTTP/2 处理逻辑
	default:
	}
	fmt.Printf("protocol: %s\n", protocl)
}

// echo 框架处理Http请求
func (v *Vortex) handleHttpWithEcho(dispatcher *Dispatcher) {

	req, err := http.ReadRequest(bufio.NewReader(dispatcher.GetReadBuffer()))
	if nil != err {
		fmt.Printf("read request error: %v\n", err)
		return
	}
	defer req.Body.Close()

	rec := httptest.NewRecorder()

	echoCtx := v.httpServ.e.NewContext(req, rec)
	v.httpServ.e.Router().Find(echoCtx.Request().Method, echoCtx.Request().URL.Path, echoCtx)
	if echoCtx.Handler() == nil {
		echoCtx.String(http.StatusNotFound, "404 Not Found")
	} else {
		if err := echoCtx.Handler(); nil != err {
			echoCtx.String(http.StatusInternalServerError, "500 Internal Server Error")
		}
	}

	resp := rec.Result()
	_ = dispatcher.Response(fmt.Appendf(nil,
		"HTTP/1.1 %d %s\r\n%s\r\n%s",
		resp.StatusCode,
		http.StatusText(resp.StatusCode),
		resp.Header,
		rec.Body.String(),
	))
}
