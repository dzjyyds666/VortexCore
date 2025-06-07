package vortex

import (
	"context"
	"io"
	"net/http"
	"time"

	vortexMw "github.com/dzjyyds666/VortexCore/middleware"
	vortexUtil "github.com/dzjyyds666/VortexCore/utils"

	"github.com/labstack/echo/v4"
)

type httpContext struct {
	echo.Context // Echo 上下文
}

func (h *httpContext) GetContext() context.Context {
	return h.Request().Context()
}

func (h *httpContext) GetEcho() echo.Context {
	return h.Context
}

type httpServer struct {
	ctx context.Context
	e   *echo.Echo // Echo 框架实例
}

func NewHttpServer(ctx context.Context, routers []*HttpRouter) *httpServer {
	e := echo.New()

	vortex := e.Group("/v1")

	for _, router := range routers {
		for _, method := range router.method {
			vortex.Add(method, router.path, func(ctx echo.Context) error {
				// 包装成自身封装的上下文
				return router.handle(&httpContext{ctx})
			}, router.ToMiddleWareList()...)
		}
	}

	return &httpServer{
		ctx: ctx,
		e:   e,
	}
}

type HttpRouter struct {
	handle      func(VortexContext) error       // 路由处理函数
	path        string                          // 路由路径
	method      []string                        // HTTP方法
	middleWares []vortexMw.VortexHttpMiddleware // 中间件
	description string                          // 路由的描述
}

// 添加 Http 路由
func AppendHttpRouter(method []string, path string, handle func(VortexContext) error, apiDescription string, middleWares ...vortexMw.VortexHttpMiddleware) *HttpRouter {
	// 中间件顺序调用 parseJwt -> 自定义中间件 -> verifyJwt
	mws := make([]vortexMw.VortexHttpMiddleware, 0)
	mws = append(mws, vortexMw.PrintRequestInfoMw(), vortexMw.PrintResponseInfoMw(), vortexMw.JwtParseMw())
	mws = append(mws, middleWares...)
	mws = append(mws, vortexMw.JwtVerifyMw())

	return &HttpRouter{
		handle:      handle,
		path:        path,
		method:      method,
		middleWares: mws,
		description: apiDescription,
	}
}

// 将 VortexHttpMiddleware 转换为 Echo 中间件列表
// 这将允许 Echo 框架使用这些中间件
func (hr *HttpRouter) ToMiddleWareList() []echo.MiddlewareFunc {
	middlewares := make([]echo.MiddlewareFunc, 0, len(hr.middleWares))
	for _, mw := range hr.middleWares {
		middlewares = append(middlewares, echo.MiddlewareFunc(mw))
	}
	return middlewares
}

type HttpOpt func(resp http.Header) http.Header

// HttpJsonResponse 返回json数据
func HttpJsonResponse(ctx VortexContext, code int, data interface{}, opts ...HttpOpt) error {
	echoCtx := ctx.GetEcho()
	// 设置响应的请求头
	for _, opt := range opts {
		opt(echoCtx.Response().Header())
	}

	return echoCtx.JSON(code, VortexHttpResponse{
		Code: code,
		Body: data,
		Info: struct {
			Url  string `json:"url,omitempty"`  // 响应的url
			Time int64  `json:"time,omitempty"` // 响应时间
		}{
			Url:  echoCtx.Request().URL.String(),
			Time: time.Now().Unix(),
		},
	})
}

// 流式返回数据
func HttpStreamResponse(ctx VortexContext, code int, stream io.Reader, opts ...HttpOpt) error {
	echoCtx := ctx.GetEcho()
	for _, opt := range opts {
		opt(echoCtx.Response().Header())
	}
	contentType := echoCtx.Response().Header().Get(vortexUtil.VortexHeaders.ContentType.S())
	if len(contentType) <= 0 {
		return echoCtx.Stream(code, "application/octet-stream", stream)
	} else {
		return echoCtx.Stream(code, contentType, stream)
	}
}
