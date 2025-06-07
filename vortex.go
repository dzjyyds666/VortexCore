package vortex

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"

	vortexUtil "github.com/dzjyyds666/VortexCore/utils"

	"github.com/dzjyyds666/opensource/logx"
)

var Transport = struct {
	TCP string
	UDP string
}{
	TCP: "tcp",
	UDP: "udp",
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
	hideBanner bool               // 是否隐藏启动旗帜，默认为false，显示banner
}

// 启动服务
func NewVortexCore(ctx context.Context, opts ...Option) *Vortex {
	vortex := &Vortex{
		ctx:       ctx,
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
		case vortexUtil.Http1:
			router := prepareDefaultHttpRouter()
			vortex.httpRouter = append(vortex.httpRouter, router...)
			vortex.httpServ = NewHttpServer(ctx, vortex.httpRouter)
		}
	}

	return vortex
}

// 开启端口监听，先判断当前请求的协议，然后选择对应的协议进行处理
func (v *Vortex) Start() {
	ln, err := net.Listen(v.transport, fmt.Sprintf(":%s", v.port))
	if nil != err {
		panic(err)
	}
	defer ln.Close()

	if !v.hideBanner {
		vortexUtil.ShowBanner(v.port)
	}

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
		err := conn.Close()
		if nil != err {
			fmt.Printf("close error: %v\n", err)
		}
	}()

	d := NewDispatcher(ctx, conn)
	protocl, err := d.Parse()
	// 关闭连接
	if nil != err || protocl == "unknown" {
		fmt.Printf("parse error: %v\n", err)
		return
	}

	switch protocl {
	case vortexUtil.Http1:
		// 使用echo框架处理 HTTP/1.1 请求
		err := v.handleHttpWithEcho(d)
		if nil != err {
			d.Response([]byte("500 Internal Server Error"))
		}
	case vortexUtil.WebSocket:
		// 使用 WebSocket 处理逻辑
	case vortexUtil.Http2:
		// 使用 HTTP/2 处理逻辑
	default:
	}
}

// echo 框架处理Http请求
func (v *Vortex) handleHttpWithEcho(dispatcher *Dispatcher) error {

	req, err := http.ReadRequest(bufio.NewReader(dispatcher.GetReadBuffer()))
	if nil != err {
		fmt.Printf("read request error: %v\n", err)
		return err
	}
	defer req.Body.Close()

	rec := httptest.NewRecorder()

	echoCtx := v.httpServ.e.NewContext(req, rec)
	v.httpServ.e.Router().Find(echoCtx.Request().Method, echoCtx.Request().URL.Path, echoCtx)
	if echoCtx.Handler() == nil {
		echoCtx.String(http.StatusNotFound, "404 Not Found")
	} else {
		if err := echoCtx.Handler()(echoCtx); nil != err {
			echoCtx.String(http.StatusInternalServerError, "500 Internal Server Error")
		}
	}

	resp := rec.Result()

	var buf bytes.Buffer
	err = resp.Write(&buf)
	if nil != err {
		return err
	}
	err = dispatcher.Response(buf.Bytes())
	return err
}

type Option func(*Vortex)

// 设置自定义日志
func WithCustomLogger(logPath string, logLevel logx.LogLevel, maxSizeMB int64, consoleOut bool) Option {
	return func(v *Vortex) {
		if err := vortexUtil.InitVortexLog(logPath, logLevel, maxSizeMB, consoleOut); nil != err {
			panic(fmt.Sprintf("init vortex log error: %v", err))
		}
	}
}

// 设置默认日志
func WithDefaultLogger() Option {
	return func(v *Vortex) {
		if err := vortexUtil.InitVortexLog("logs/stdout.log", logx.DEBUG, 10, true); nil != err {
			panic(fmt.Sprintf("init vortex log error: %v", err))
		}
	}
}

// 设置监听端口
func WithListenPort(port string) Option {
	return func(v *Vortex) {
		v.port = port
	}
}

// 是否显示启动旗帜
func WithHideBanner(show bool) Option {
	return func(v *Vortex) {
		v.hideBanner = show
	}
}

// 设置传输协议
func WithTransport(transport string) Option {
	return func(v *Vortex) {
		v.transport = transport
	}
}

// 设置支持的协议
func WithProtocol(protocols ...string) Option {
	return func(v *Vortex) {
		if v.protocol == nil {
			v.protocol = make([]string, 0)
		}
		v.protocol = append(v.protocol, protocols...)
	}
}

// 设置自定义Http路由
func WithHttpRouter(routers []*HttpRouter) Option {
	return func(v *Vortex) {
		if v.httpRouter == nil {
			v.httpRouter = make([]*HttpRouter, 0)
		}
		v.httpRouter = append(v.httpRouter, routers...)
	}
}
